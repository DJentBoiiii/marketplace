package filetransfer

import (
	"fmt"
	"os"

	"github.com/DjentBoiiii/marketplace/internal/auth"
	"github.com/DjentBoiiii/marketplace/internal/render"
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

func SetupUploadHandlers(app *fiber.App) {
	app.Get("/upload", auth.LoginRequired(), func(c *fiber.Ctx) error {
		return render.RenderTemplate(c, "upload.html")
	})
	app.Post("/upload", auth.LoginRequired(), UploadFile)
	// app.Get("/download/:filePath", DownloadFile)
	app.Get("/delete", auth.LoginRequired(), func(c *fiber.Ctx) error {
		return render.RenderTemplate(c, "delete.html")
	})
	app.Post("/delete", auth.LoginRequired(), DeleteFile)
	app.Post("/download", auth.LoginRequired(), DownloadFile)
}
