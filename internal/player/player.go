package player

import (
	"context"
	"database/sql"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

var minioClient *minio.Client
var bucketName = "djent"
var DB_USER = os.Getenv("MYSQL_USER")
var DB_PASSWORD = os.Getenv("MYSQL_PASSWORD")
var DB_NAME = os.Getenv("MYSQL_DATABASE")
var _ = godotenv.Load("/marketplace/.env")

func init() {
	endpoint := os.Getenv("MINIO_ENDPOINT")
	accessKeyID := os.Getenv("MINIO_ACCESS_KEY")
	secretAccessKey := os.Getenv("MINIO_SECRET_KEY")
	useSSL := false

	var err error
	minioClient, err = minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		log.Fatalf("Failed to initialize MinIO client: %v", err)
	}
}

func GetAudio(c *fiber.Ctx) error {
	trackID := c.Params("track_id")

	db, err := sql.Open("mysql", DB_USER+":"+DB_PASSWORD+"@tcp(boku-no-sukele:3306)/"+DB_NAME)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to connect to database"})
	}
	defer db.Close()

	var path string
	err = db.QueryRow(`SELECT CONCAT(vendor, '/', type, '/', name) FROM Products WHERE id = ?`, trackID).Scan(&path)
	if err != nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": "Audio file not found in DB"})
	}

	object, err := minioClient.GetObject(context.Background(), bucketName, path, minio.GetObjectOptions{})
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch audio file"})
	}
	defer object.Close()

	stat, err := object.Stat()
	if err != nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": "Audio file not found in MinIO"})
	}

	c.Set("Content-Type", "audio/mpeg")
	c.Set("Content-Length", strconv.FormatInt(stat.Size, 10))

	_, err = io.Copy(c.Response().BodyWriter(), object)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to stream audio file"})
	}

	return nil
}

func RegisterRoutes(app *fiber.App) {
	app.Get("/audio/:track_id", GetAudio)
}
