// internal/playlist/viewPlaylist.go
package playlist

import (
	"github.com/DjentBoiiii/marketplace/internal/auth"
	"github.com/DjentBoiiii/marketplace/internal/render"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gofiber/fiber/v2"
)

func ViewPlaylist(c *fiber.Ctx) error {
	user, err := auth.GetUserData(c)
	if err != nil {
		return c.Status(500).SendString("Помилка отримання даних користувача")
	}

	if user.Username == "" {
		return c.Redirect("/login")
	}

	playlistID := c.Params("id")
	if playlistID == "" {
		return c.Status(400).SendString("Не вказано ID плейлиста")
	}

	// Перевірка чи належить плейлист цьому користувачу
	var playlistUserID int
	var playlistName string
	err = DB.QueryRow("SELECT user_id, name FROM Playlists WHERE id = ?", playlistID).Scan(&playlistUserID, &playlistName)
	if err != nil {
		return c.Status(404).SendString("Плейлист не знайдено")
	}

	if playlistUserID != user.Id {
		return c.Status(403).SendString("Ви не маєте доступу до цього плейлиста")
	}

	playlistItems, err := GetPlaylistItems(playlistID)
	if err != nil {
		return c.Status(500).SendString("Помилка отримання вмісту плейлиста: " + err.Error())
	}

	data := render.TemplateData{
		"playlist_id":   playlistID,
		"playlist_name": playlistName,
		"items":         playlistItems,
		"user":          user,
	}

	return render.RenderTemplate(c, "playlist_detail.html", data)
}
