package filetransfer

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"

	"github.com/DjentBoiiii/marketplace/internal"
	handlers "github.com/DjentBoiiii/marketplace/internal/handlers/auth"
	"github.com/gofiber/fiber/v2"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

var (
	MinioEndpoint  = os.Getenv("MINIO_ENDPOINT")
	MinioAccessKey = os.Getenv("MINIO_ACCESS_KEY")
	MinioSecretKey = os.Getenv("MINIO_SECRET_KEY")
	MinioClient    *minio.Client
)

func init() {
	var err error
	MinioClient, err = minio.New(MinioEndpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(MinioAccessKey, MinioSecretKey, ""),
		Secure: false,
	})
	if err != nil {
		fmt.Println(err)
		fmt.Println(MinioClient)
		panic(fmt.Sprintf("Не вдалося підключитися до MinIO: %v", err))
	}
}

func UploadFile(c *fiber.Ctx) error {
	user, err := handlers.GetUserData(c)
	if err != nil {
		return c.Status(500).SendString("Помилка отримання даних користувача")
	}

	file, err := c.FormFile("file")
	if err != nil {
		return c.Status(400).SendString("Файл не надано")
	}

	category := c.FormValue("category")
	if category != "audio" && category != "midi" && category != "samples" {
		return c.Status(400).SendString("Некоректна категорія")
	}

	src, err := file.Open()
	if err != nil {
		return c.Status(500).SendString("Помилка відкриття файлу")
	}
	defer src.Close()

	if category == "audio" {
		if !isAudio(file.Filename) {
			return c.Status(400).SendString("Файл не є аудіоформатом")
		}
	} else {
		if !isArchive(file.Filename) {
			return c.Status(400).SendString("Файл має бути архівом")
		}
		if hasExecutableFiles(src) {
			return c.Status(400).SendString("Архів містить виконувані файли")
		}
	}

	bucketName := strings.ToLower(user.Username)
	objectPath := fmt.Sprintf("%s/%s/%s", user.Username, category, file.Filename)

	exists, _ := MinioClient.BucketExists(c.Context(), bucketName)
	if !exists {
		err := MinioClient.MakeBucket(c.Context(), bucketName, minio.MakeBucketOptions{})
		if err != nil {
			return c.Status(500).SendString("Не вдалося створити бакет")
		}
	}

	_, err = MinioClient.PutObject(c.Context(), bucketName, objectPath, src, file.Size, minio.PutObjectOptions{ContentType: file.Header["Content-Type"][0]})
	if err != nil {
		return c.Status(500).SendString("Помилка завантаження файлу")
	}

	return c.SendString("Файл успішно завантажено")
}

func isAudio(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	audioFormats := []string{".mp3", ".wav", ".flac", ".ogg", ".m4a"}
	for _, format := range audioFormats {
		if ext == format {
			return true
		}
	}
	return false
}
func isArchive(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	return ext == ".zip"
}

func hasExecutableFiles(file multipart.File) bool {
	buf := new(bytes.Buffer)
	_, _ = io.Copy(buf, file)
	r, err := zip.NewReader(bytes.NewReader(buf.Bytes()), int64(buf.Len()))
	if err != nil {
		return false
	}
	for _, f := range r.File {
		if strings.HasSuffix(f.Name, ".exe") || strings.HasSuffix(f.Name, ".bat") || strings.HasSuffix(f.Name, ".sh") {
			return true
		}
	}
	return false
}

func SetupUploadHandlers(app *fiber.App) {
	app.Get("/upload", handlers.LoginRequired(), func(c *fiber.Ctx) error {
		return internal.RenderTemplate(c, "upload.html")
	})
	app.Post("/upload", handlers.LoginRequired(), UploadFile)
}
