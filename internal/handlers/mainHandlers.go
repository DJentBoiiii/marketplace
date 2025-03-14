package handlers

import (
	"github.com/DjentBoiiii/marketplace/internal"
	"github.com/DjentBoiiii/marketplace/internal/filetransfer"
	handlers "github.com/DjentBoiiii/marketplace/internal/handlers/auth"
	"github.com/gofiber/fiber/v2"
)

func index(c *fiber.Ctx) error {
	user, err := handlers.GetUserData(c)
	if err != nil {
		return err
	}
	return internal.RenderTemplate(c, "index.html", [2]interface{}{"user", user})
}

func SetupHandlers(app *fiber.App) {
	app.Get("/", index)
	handlers.SetupAuthHandlers(app)
	filetransfer.SetupUploadHandlers(app)

}
