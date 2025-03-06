package handlers

import (
	"github.com/gofiber/fiber/v2"
)

const HTML_PATH = "/marketplace/web/static/templates"

func index(c *fiber.Ctx) error {
	return c.SendFile(HTML_PATH + "/index.html")
}

func SetupHandlers(app *fiber.App) {
	app.Get("/", index)
}
