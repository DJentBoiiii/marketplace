package filetransfer

import (
	"database/sql"
	"fmt"
	"io"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/minio/minio-go/v7"
)

func DownloadFile(c *fiber.Ctx) error {

	audio_id := c.Params("id")

	db, _ := sql.Open("mysql", DB_USER+":"+DB_PASSWORD+"@tcp(boku-no-sukele:3306)/"+DB_NAME)
	defer db.Close()
	var vendor, name, fileType string
	err := db.QueryRow("SELECT vendor, type, name FROM Products WHERE id = ?", audio_id).Scan(&vendor, &fileType, &name)
	if err != nil {
		return c.Status(500).SendString("Помилка отримання даних продукту")
	}

	bucketName := strings.ToLower("djent")
	objectPath := fmt.Sprintf("%s/%s/%s", vendor, fileType, name)

	object, err := MinioClient.GetObject(c.Context(), bucketName, objectPath, minio.GetObjectOptions{})
	if err != nil {
		return c.Status(404).SendString("Файл не знайдено на сервері")
	}

	c.Set("Content-Type", "application/octet-stream")
	c.Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", name))
	_, err = io.Copy(c.Response().BodyWriter(), object)
	if err != nil {
		return c.Status(500).SendString("Помилка завантаження файлу")
	}

	return nil
}
