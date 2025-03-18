// internal/playlist/removeFromPlaylist.go
package playlist

import (
	"github.com/DjentBoiiii/marketplace/internal/auth"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gofiber/fiber/v2"
)

func RemoveFromPlaylist(c *fiber.Ctx) error {
	user, err := auth.GetUserData(c)
	if err != nil {
		return c.Status(500).SendString("Помилка отримання даних користувача")
	}

	if user.Username == "" {
		return c.Redirect("/login")
	}

	itemID := c.FormValue("item_id")
	playlistID := c.FormValue("playlist_id")

	if itemID == "" || playlistID == "" {
		return c.Status(400).SendString("Не вказано ID елементу або плейлиста")
	}

	// Перевірка чи належить плейлист цьому користувачу
	var playlistUserID int
	err = DB.QueryRow("SELECT user_id FROM Playlists WHERE id = ?", playlistID).Scan(&playlistUserID)
	if err != nil {
		return c.Status(500).SendString("Помилка перевірки плейлиста")
	}

	if playlistUserID != user.Id {
		return c.Status(403).SendString("Ви не маєте доступу до цього плейлиста")
	}

	// Видалення елементу з плейлиста
	_, err = DB.Exec("DELETE FROM PlaylistItems WHERE id = ? AND playlist_id = ?",
		itemID, playlistID)
	if err != nil {
		return c.Status(500).SendString("Помилка видалення з плейлиста")
	}

	return c.Redirect("/playlist/" + playlistID)
}
