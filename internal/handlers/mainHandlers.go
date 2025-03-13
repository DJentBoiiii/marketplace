package handlers

import (
	"github.com/DjentBoiiii/marketplace/internal"
	handlers "github.com/DjentBoiiii/marketplace/internal/handlers/auth"
	"github.com/gofiber/fiber/v2"
)

func index(c *fiber.Ctx) error {
	return internal.RenderTemplate(c, "index.html", nil)
}

func SetupHandlers(app *fiber.App) {
	app.Get("/", index)
	handlers.SetupAuthHandlers(app)
}
