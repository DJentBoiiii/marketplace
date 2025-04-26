package filetransfer

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/minio/minio-go/v7"
)

func DeleteFile(c *fiber.Ctx) error {

	filePath := "/DJent/audio/Kira.mp3"

	decodedFilePath, err := url.QueryUnescape(filePath)
	if err != nil {
		return c.Status(400).SendString("Невірний шлях до файлу")
	}
	bucketName := strings.ToLower("dyploma-marketplace-products")
	exists, err := MinioClient.BucketExists(context.Background(), bucketName)
	if err != nil {
		return c.Status(500).SendString("Помилка перевірки бакету: " + err.Error())
	}
	if !exists {
		return c.Status(404).SendString("Бакет не знайдений")
	}
	err = MinioClient.RemoveObject(context.Background(), bucketName, decodedFilePath, minio.RemoveObjectOptions{})
	if err != nil {
		return c.Status(500).SendString("Помилка видалення файлу: " + err.Error())
	}
	return c.SendString(fmt.Sprintf("Файл %s успішно видалено", decodedFilePath))
}
