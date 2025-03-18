// internal/playlist/deletePlaylist.go
package playlist

import (
	"github.com/DjentBoiiii/marketplace/internal/auth"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gofiber/fiber/v2"
)

func DeletePlaylist(c *fiber.Ctx) error {
	user, err := auth.GetUserData(c)
	if err != nil {
		return c.Status(500).SendString("Помилка отримання даних користувача")
	}

	if user.Username == "" {
		return c.Redirect("/login")
	}

	playlistID := c.FormValue("playlist_id")
	if playlistID == "" {
		return c.Status(400).SendString("Не вказано ID плейлиста")
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

	// Видалення плейлиста та всіх його елементів
	_, err = DB.Exec("DELETE FROM Playlists WHERE id = ?", playlistID)
	if err != nil {
		return c.Status(500).SendString("Помилка видалення плейлиста")
	}

	return c.Redirect("/playlists")
}
