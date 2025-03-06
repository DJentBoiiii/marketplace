package service

import (
	"fmt"

	"github.com/DjentBoiiii/marketplace/internal/handlers"
	"github.com/gofiber/fiber/v2"
)

func StartService() {
	app := fiber.New()
	app.Static("/static", "../../web/static")
	fmt.Println("Service starts on port 69420")
	handlers.SetupHandlers(app)
	app.Listen(":3000")
}
