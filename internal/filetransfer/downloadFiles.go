package filetransfer

import (
	"fmt"
	"io"
	"strings"

	"github.com/DjentBoiiii/marketplace/internal/auth"
	"github.com/gofiber/fiber/v2"
	"github.com/minio/minio-go/v7"
)

func DownloadFile(c *fiber.Ctx) error {

	user, err := auth.GetUserData(c)
	if err != nil {
		return c.Status(500).SendString("Помилка отримання даних користувача")
	}

	filePath := "samples/sample_test.zip"
	if filePath == "" {
		return c.Status(400).SendString("Файл не знайдено")
	}

	bucketName := strings.ToLower(user.Username)
	objectPath := fmt.Sprintf("%s/%s", user.Username, filePath)

	object, err := MinioClient.GetObject(c.Context(), bucketName, objectPath, minio.GetObjectOptions{})
	if err != nil {
		return c.Status(404).SendString("Файл не знайдено на сервері")
	}

	c.Set("Content-Type", "application/octet-stream")
	c.Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filePath))
	_, err = io.Copy(c.Response().BodyWriter(), object)
	if err != nil {
		return c.Status(500).SendString("Помилка завантаження файлу")
	}

	return nil
}
