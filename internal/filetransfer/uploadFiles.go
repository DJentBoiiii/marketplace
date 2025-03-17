package filetransfer

import (
	"bytes"
	"fmt"
	"io"
	"strings"

	"github.com/DjentBoiiii/marketplace/internal/auth"
	"github.com/gofiber/fiber/v2"
	"github.com/minio/minio-go/v7"
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

	category := c.FormValue("category")
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

	reader := bytes.NewReader(fileData.Bytes())

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
		reader = bytes.NewReader(fileData.Bytes())
	}

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

	_, err = MinioClient.PutObject(c.Context(), bucketName, objectPath, reader, int64(fileData.Len()), minio.PutObjectOptions{ContentType: file.Header["Content-Type"][0]})
	if err != nil {
		fmt.Println(err)
		return c.Status(500).SendString("Помилка завантаження файлу")
	}

	return c.SendString("Файл успішно завантажено")
}
