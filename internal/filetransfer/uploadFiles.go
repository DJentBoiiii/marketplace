package filetransfer

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/DjentBoiiii/marketplace/internal/auth"
	"github.com/minio/minio-go/v7"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gofiber/fiber/v2"
)

func UploadFile(c *fiber.Ctx) error {
	user, err := auth.GetUserData(c)
	if err != nil {
		return c.Status(500).SendString("Помилка отримання даних користувача")
	}

	file, err := c.FormFile("file")
	if err != nil {
		return c.Status(400).SendString("Файл не надано")
	}

	image, err := c.FormFile("image")
	if err != nil {
		return c.Status(400).SendString("Зображення не надано")
	}

	typeVal := c.FormValue("type")
	description := c.FormValue("description")
	price, _ := strconv.Atoi(c.FormValue("price"))

	// Add genre handling
	var genre string
	if typeVal == "audio" {
		genre = c.FormValue("genre")
	} else if typeVal == "midi" || typeVal == "samples" {
		genre = c.FormValue("subtype")
	}

	src, err := file.Open()
	if err != nil {
		return c.Status(500).SendString("Помилка відкриття файлу")
	}
	defer src.Close()

	fileData := new(bytes.Buffer)
	if _, err := io.Copy(fileData, src); err != nil {
		return c.Status(500).SendString("Помилка читання файлу")
	}

	if typeVal == "audio" && !isAudio(file.Filename) {
		return c.Status(400).SendString("Файл має бути аудіоформатом")
	}
	if (typeVal == "midi" || typeVal == "samples") && !isArchive(file.Filename) {
		return c.Status(400).SendString("Файл має бути ZIP-архівом")
	}

	imgSrc, err := image.Open()
	if err != nil {
		return c.Status(500).SendString("Помилка відкриття зображення")
	}
	defer imgSrc.Close()

	imgData := new(bytes.Buffer)
	if _, err := io.Copy(imgData, imgSrc); err != nil {
		return c.Status(500).SendString("Помилка читання зображення")
	}

	imgDir := "/marketplace/web/static/uploads/products/" + typeVal + "/"

	imgPath := fmt.Sprintf("%s%s", imgDir, image.Filename)
	imgDBPath := "static/uploads/products/" + typeVal + "/" + image.Filename
	imgFile, err := os.Create(imgPath)
	if err != nil {
		fmt.Println(err)
		return c.Status(500).SendString("Не вдалося створити файл зображення")
	}
	defer imgFile.Close()

	if _, err := imgFile.Write(imgData.Bytes()); err != nil {
		return c.Status(500).SendString("Помилка запису зображення")
	}

	bucketName := strings.ToLower("dyploma-marketplace-products")
	objectPath := fmt.Sprintf("%s/%s/%s", user.Username, typeVal, file.Filename)

	exists, _ := MinioClient.BucketExists(c.Context(), bucketName)
	if !exists {
		err := MinioClient.MakeBucket(c.Context(), bucketName, minio.MakeBucketOptions{})
		if err != nil {
			return c.Status(500).SendString("Не вдалося створити бакет")
		}
	}

	_, err = MinioClient.PutObject(
		c.Context(),
		bucketName,
		objectPath,
		bytes.NewReader(fileData.Bytes()),
		int64(fileData.Len()),
		minio.PutObjectOptions{ContentType: file.Header["Content-Type"][0]},
	)
	if err != nil {
		return c.Status(500).SendString("Помилка завантаження файлу")
	}
	filename := strings.TrimSuffix(file.Filename, filepath.Ext(file.Filename))
	extension := filepath.Ext(file.Filename)

	_, err = DB.Exec(
		"INSERT INTO Products (name, type, price, description, vendor, image_url, Extension, genre) VALUES (?, ?, ?, ?, ?, ?, ?, ?)",
		filename, typeVal, price, description, user.Username, imgDBPath, extension, genre,
	)
	if err != nil {
		fmt.Println(err)
		return c.Status(500).SendString("Помилка запису в базу")
	}

	// Set is_artist flag to true for the user who uploaded a product
	_, err = DB.Exec("UPDATE Users SET is_artist = TRUE WHERE username = ?", user.Username)
	if err != nil {
		fmt.Println("Error updating user as artist:", err)
		// Continue execution even if setting the artist flag fails
	}

	err = sendEmbedRequest(objectPath, "(C)"+user.Username)
	if err != nil {
		fmt.Println("Помилка вмонтування водяного знаку:", err)
		return c.Status(500).SendString("Помилка вмонтування водяного знаку")
	}

	fmt.Println("File uploaded successfully:", filename)

	return c.SendString("Файл успішно завантажено")
}

func sendEmbedRequest(s3Key, message string) error {
	var b bytes.Buffer
	w := multipart.NewWriter(&b)

	// Додаємо s3_key
	err := w.WriteField("s3_key", s3Key)
	if err != nil {
		return err
	}

	// Додаємо message
	err = w.WriteField("message", message)
	if err != nil {
		return err
	}

	w.Close()

	req, err := http.NewRequest("POST", "http://fastapi-app:8000/embed", &b)
	if err != nil {
		return err
	}

	// Встановлюємо правильний Content-Type
	req.Header.Set("Content-Type", w.FormDataContentType())

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("FastAPI повернув помилку %s: %s", resp.Status, string(bodyBytes))
	}

	return nil
}
