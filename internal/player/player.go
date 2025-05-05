package player

import (
	"context"
	"database/sql"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"

	_ "github.com/go-sql-driver/mysql" // Import MySQL driver
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

var minioClient *minio.Client
var bucketName = "dyploma-marketplace-products"
var DB_USER string
var DB_PASSWORD string
var DB_NAME string
var DB_HOST string

func init() {
	// Load environment variables first
	err := godotenv.Load("/marketplace/.env")
	if err != nil {
		log.Printf("Warning: Error loading .env file: %v", err)
	}

	// Set database connection variables
	DB_USER = os.Getenv("MYSQL_USER")
	DB_PASSWORD = os.Getenv("MYSQL_PASSWORD")
	DB_NAME = os.Getenv("MYSQL_DATABASE")
	DB_HOST = os.Getenv("DB_HOST")

	if DB_HOST == "" {
		DB_HOST = "localhost"
		log.Printf("DB_HOST not set, defaulting to %s", DB_HOST)
	}

	// Initialize MinIO client
	endpoint := os.Getenv("MINIO_ENDPOINT")
	accessKeyID := os.Getenv("MINIO_ACCESS_KEY")
	secretAccessKey := os.Getenv("MINIO_SECRET_KEY")
	useSSL := false

	minioClient, err = minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		log.Fatalf("Failed to initialize MinIO client: %v", err)
	}

	log.Printf("Audio player initialized with endpoint: %s", endpoint)
}

func GetAudio(c *fiber.Ctx) error {
	trackID := c.Params("track_id")
	log.Printf("GetAudio called with track_id: %s", trackID)

	db, err := sql.Open("mysql", DB_USER+":"+DB_PASSWORD+"@tcp("+DB_HOST+":3306)/"+DB_NAME)
	if err != nil {
		log.Printf("Database connection error: %v", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to connect to database"})
	}
	defer db.Close()
	log.Printf("Database connected successfully")

	var path string
	err = db.QueryRow(`SELECT CONCAT(vendor, '/watermarked/', name, '.wav') FROM Products WHERE id = ?`, trackID).Scan(&path)
	if err != nil {
		log.Printf("Database query error for track_id %s: %v", trackID, err)
		return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": "Audio file not found in DB"})
	}
	log.Printf("Found audio file path in DB: %s", path)

	object, err := minioClient.GetObject(context.Background(), bucketName, path, minio.GetObjectOptions{})
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch audio file"})
	}
	defer object.Close()

	stat, err := object.Stat()
	if err != nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": "Audio file not found in MinIO"})
	}

	contentType := stat.ContentType
	if contentType == "" {
		contentType = "application/octet-stream"
	}
	c.Set("Content-Type", contentType)
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
