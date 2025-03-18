// internal/playlist/addToPlaylist.go
package playlist

import (
	"github.com/DjentBoiiii/marketplace/internal/auth"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gofiber/fiber/v2"
)

func AddToPlaylist(c *fiber.Ctx) error {
	user, err := auth.GetUserData(c)
	if err != nil {
		return c.Status(500).SendString("Помилка отримання даних користувача")
	}

	if user.Username == "" {
		return c.Redirect("/login")
	}

	productID := c.FormValue("product_id")
	playlistID := c.FormValue("playlist_id")

	if productID == "" || playlistID == "" {
		return c.Status(400).SendString("Не вказано ID продукту або плейлиста")
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

	// Перевірка чи продукт вже є в плейлисті
	var count int
	err = DB.QueryRow("SELECT COUNT(*) FROM PlaylistItems WHERE playlist_id = ? AND product_id = ?",
		playlistID, productID).Scan(&count)
	if err != nil {
		return c.Status(500).SendString("Помилка перевірки елементу плейлиста")
	}

	if count > 0 {
		return c.Status(400).SendString("Цей продукт вже є у вашому плейлисті")
	}

	// Додавання продукту до плейлиста
	_, err = DB.Exec("INSERT INTO PlaylistItems (playlist_id, product_id) VALUES (?, ?)",
		playlistID, productID)
	if err != nil {
		return c.Status(500).SendString("Помилка додавання до плейлиста")
	}

	referer := c.Get("Referer")
	if referer == "" {
		referer = "/playlist/" + playlistID
	}
	return c.Redirect(referer)
}
