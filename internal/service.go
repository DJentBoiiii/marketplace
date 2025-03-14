package internal

import (
	"fmt"

	"github.com/DjentBoiiii/marketplace/internal/auth"
	"github.com/DjentBoiiii/marketplace/internal/filetransfer"
	"github.com/gofiber/fiber/v2"
)

func SetupHandlers(app *fiber.App) {
	app.Get("/", index)
	auth.SetupAuthHandlers(app)
	filetransfer.SetupUploadHandlers(app)

}

func StartService() {
	app := fiber.New()
	app.Static("/static", "../../web/static")
	fmt.Println("Service starts on port 69420")
	SetupHandlers(app)
	app.Listen(":3000")
}
