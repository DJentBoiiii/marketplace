package auth

import (
	"github.com/DjentBoiiii/marketplace/internal/render"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gofiber/fiber/v2"
)

func Profile(c *fiber.Ctx) error {
	username := c.Params("username")
	user, err := GetUserFromDB(username)
	if err != nil {
		return c.Status(404).SendString("Користувач не знайдений")
	}

	return render.RenderTemplate(c, "profile.html",
		[2]interface{}{"user", user},
	)
}
