package filetransfer

import (
	"bytes"
	"database/sql"
	"fmt"
	"io"
	"os"
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

	category := c.FormValue("category")
	description := c.FormValue("description")
	price, _ := strconv.Atoi(c.FormValue("price"))

	if category != "audio" && category != "midi" && category != "samples" {
		return c.Status(400).SendString("Некоректна категорія")
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

	imgSrc, err := image.Open()
	if err != nil {
		return c.Status(500).SendString("Помилка відкриття зображення")
	}
	defer imgSrc.Close()

	imgData := new(bytes.Buffer)
	if _, err := io.Copy(imgData, imgSrc); err != nil {
		return c.Status(500).SendString("Помилка читання зображення")
	}

	// Збереження зображення в локальну директорію
	imgDir := "static/images/"
	if err := os.MkdirAll(imgDir, os.ModePerm); err != nil {
		return c.Status(500).SendString("Не вдалося створити директорію для зображень")
	}

	imgPath := fmt.Sprintf("%s%s", imgDir, image.Filename)
	imgFile, err := os.Create(imgPath)
	if err != nil {
		return c.Status(500).SendString("Не вдалося створити файл зображення")
	}
	defer imgFile.Close()

	if _, err := imgFile.Write(imgData.Bytes()); err != nil {
		return c.Status(500).SendString("Помилка запису зображення")
	}

	// Перевірка категорії файлів
	switch category {
	case "audio":
		if !isAudio(file.Filename) {
			return c.Status(400).SendString("Файл не є аудіоформатом")
		}
	case "midi":
		if !isArchive(file.Filename) {
			return c.Status(400).SendString("Файл має бути архівом")
		}
		if !hasMIDI(fileData.Bytes()) {
			return c.Status(400).SendString("Архів не містить MIDI-файлів")
		}
	case "samples":
		if !isArchive(file.Filename) {
			return c.Status(400).SendString("Файл має бути архівом")
		}
		if hasExecutableFiles(src) {
			return c.Status(400).SendString("Архів містить виконувані файли")
		}
	}

	// Завантаження файлу на MinIO
	bucketName := strings.ToLower(user.Username)
	objectPath := fmt.Sprintf("%s/%s/%s", user.Username, category, file.Filename)

	exists, _ := MinioClient.BucketExists(c.Context(), bucketName)
	if !exists {
		err := MinioClient.MakeBucket(c.Context(), bucketName, minio.MakeBucketOptions{})
		if err != nil {
			fmt.Println(err)
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
		fmt.Println(err)
		return c.Status(500).SendString("Помилка завантаження файлу")
	}

	// Збереження даних у базу даних
	db, _ := sql.Open("mysql", DB_USER+":"+DB_PASSWORD+"@tcp(boku-no-sukele:3306)/"+DB_NAME)
	_, err = db.Exec(
		"INSERT INTO Products (name, type, price, description, vendor, product_path, product_img) VALUES (?, ?, ?, ?, ?, ?, ?)",
		file.Filename, category, price, description, user.Id, objectPath, imgPath,
	)
	if err != nil {
		fmt.Println(err)
		return c.Status(500).SendString("Помилка запису в базу")
	}

	return c.SendString("Файл успішно завантажено")
}
