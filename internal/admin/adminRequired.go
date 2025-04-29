package admin

import (
	"github.com/DjentBoiiii/marketplace/internal/auth"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gofiber/fiber/v2"
)

// AdminRequired is a middleware that restricts access to admin-only routes
// Returns a 404 error for non-admin users to hide the existence of admin pages
func AdminRequired() fiber.Handler {
	return func(c *fiber.Ctx) error {
		user, err := auth.GetUserData(c)
		if err != nil || !user.Is_admin {
			return c.Status(404).SendString("Сторінку не знайдено")
		}
		return c.Next()
	}
}
