// internal/playlist/createPlaylist.go
package playlist

import (
	"github.com/DjentBoiiii/marketplace/internal/auth"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gofiber/fiber/v2"
)

func CreatePlaylist(c *fiber.Ctx) error {
	user, err := auth.GetUserData(c)
	if err != nil {
		return c.Status(500).SendString("Помилка отримання даних користувача")
	}

	if user.Username == "" {
		return c.Redirect("/login")
	}

	playlistName := c.FormValue("playlist_name")
	if playlistName == "" {
		return c.Status(400).SendString("Назва плейлиста не може бути порожньою")
	}

	// Перевірка чи вже існує плейлист з такою назвою у цього користувача
	var count int
	err = DB.QueryRow("SELECT COUNT(*) FROM Playlists WHERE name = ? AND user_id = ?",
		playlistName, user.Id).Scan(&count)
	if err != nil {
		return c.Status(500).SendString("Помилка перевірки плейлиста")
	}

	if count > 0 {
		return c.Status(400).SendString("Плейлист з такою назвою вже існує")
	}

	// Створення нового плейлиста
	_, err = DB.Exec("INSERT INTO Playlists (name, user_id) VALUES (?, ?)",
		playlistName, user.Id)
	if err != nil {
		return c.Status(500).SendString("Помилка створення плейлиста")
	}

	referer := c.Get("Referer")
	if referer == "" {
		referer = "/playlists"
	}
	return c.Redirect(referer)
}
