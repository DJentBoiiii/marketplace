// internal/playlist/viewPlaylists.go
package playlist

import (
	"github.com/DjentBoiiii/marketplace/internal/auth"
	"github.com/DjentBoiiii/marketplace/internal/render"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gofiber/fiber/v2"
)

func ViewPlaylists(c *fiber.Ctx) error {
	user, err := auth.GetUserData(c)
	if err != nil {
		return c.Status(500).SendString("Помилка отримання даних користувача")
	}

	if user.Username == "" {
		return c.Redirect("/login")
	}

	playlists, err := GetUserPlaylists(user.Id)
	if err != nil {
		return c.Status(500).SendString("Помилка отримання плейлистів: " + err.Error())
	}

	return render.RenderTemplate(c, "playlists.html",
		[2]interface{}{"playlists", playlists},
		[2]interface{}{"user", user},
	)
}
