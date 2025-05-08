package auth

import (
	"github.com/DjentBoiiii/marketplace/internal/render"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gofiber/fiber/v2"
)

func Profile(c *fiber.Ctx) error {
	username := c.Params("username")
	user, err := GetUserWithProfile(username)
	if err != nil {
		return c.Status(404).SendString("Користувач не знайдений")
	}

	// Get current logged in user to determine if viewing own profile
	currentUser, _ := GetUserData(c)
	isOwnProfile := currentUser.Username == user.Username

	data := render.TemplateData{
		"user":         user,
		"isOwnProfile": isOwnProfile,
	}

	return render.RenderTemplate(c, "profile.html", data)
}
