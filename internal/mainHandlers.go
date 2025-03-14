package internal

import (
	"github.com/DjentBoiiii/marketplace/internal/auth"
	"github.com/DjentBoiiii/marketplace/internal/render"
	"github.com/gofiber/fiber/v2"
)

func index(c *fiber.Ctx) error {
	user, err := auth.GetUserData(c)
	if err != nil {
		return err
	}
	return render.RenderTemplate(c, "index.html", [2]interface{}{"user", user})
}
