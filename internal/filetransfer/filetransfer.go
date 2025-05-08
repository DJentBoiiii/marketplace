package filetransfer

import (
	"fmt"

	"github.com/DjentBoiiii/marketplace/config"
	"github.com/DjentBoiiii/marketplace/internal/auth"
	"github.com/DjentBoiiii/marketplace/internal/db"
	"github.com/DjentBoiiii/marketplace/internal/render"
	"github.com/gofiber/fiber/v2"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

var (
	DB          = db.DB
	MinioClient = config.MinioClient
)

func init() {
	var err error
	MinioClient, err = minio.New(config.MinioEndpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(config.MinioAccessKey, config.MinioSecretKey, ""),
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
		return render.RenderTemplate(c, "upload.html", nil)
	})
	app.Post("/upload", auth.LoginRequired(), UploadFile)
	app.Get("/delete", auth.LoginRequired(), func(c *fiber.Ctx) error {
		return render.RenderTemplate(c, "delete.html", nil)
	})
	app.Post("/delete", auth.LoginRequired(), DeleteFile)
	app.Post("/download/:id", auth.LoginRequired(), DownloadFile)
}
