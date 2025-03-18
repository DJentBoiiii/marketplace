package filetransfer

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/DjentBoiiii/marketplace/internal/auth"
	"github.com/DjentBoiiii/marketplace/internal/render"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

var (
	MinioEndpoint  = os.Getenv("MINIO_ENDPOINT")
	MinioAccessKey = os.Getenv("MINIO_ACCESS_KEY")
	MinioSecretKey = os.Getenv("MINIO_SECRET_KEY")
	MinioClient    *minio.Client
)

var (
	_           = godotenv.Load("/marketplace/.env")
	DB_USER     = os.Getenv("MYSQL_USER")
	DB_PASSWORD = os.Getenv("MYSQL_PASSWORD")
	DB_NAME     = os.Getenv("MYSQL_DATABASE")
	JWT_SECRET  = os.Getenv("JWT_SECRET")
	SHA_SECRET  = os.Getenv("SHA_SECRET")
	DB          *sql.DB
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

func SetupUploadHandlers(app *fiber.App) {
	app.Get("/upload", auth.LoginRequired(), func(c *fiber.Ctx) error {
		return render.RenderTemplate(c, "upload.html")
	})
	app.Post("/upload", auth.LoginRequired(), UploadFile)
	app.Get("/delete", auth.LoginRequired(), func(c *fiber.Ctx) error {
		return render.RenderTemplate(c, "delete.html")
	})
	app.Post("/delete", auth.LoginRequired(), DeleteFile)
	app.Post("/download", auth.LoginRequired(), DownloadFile)
}
