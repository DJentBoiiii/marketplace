// internal/playlist/createAndAddToPlaylist.go
package playlist

import (
	"fmt"

	"github.com/DjentBoiiii/marketplace/internal/auth"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gofiber/fiber/v2"
)

func CreateAndAddToPlaylist(c *fiber.Ctx) error {
	user, err := auth.GetUserData(c)
	if err != nil {
		return c.Status(500).SendString("Помилка отримання даних користувача")
	}

	if user.Username == "" {
		return c.Redirect("/login")
	}

	productID := c.FormValue("product_id")
	playlistName := c.FormValue("playlist_name")

	if productID == "" || playlistName == "" {
		return c.Status(400).SendString("Не вказано ID продукту або назву плейлиста")
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

	// Використовуємо транзакцію для створення плейлиста і додавання продукту
	tx, err := DB.Begin()
	if err != nil {
		return c.Status(500).SendString("Помилка створення транзакції")
	}
	defer tx.Rollback()

	// Створення нового плейлиста
	result, err := tx.Exec("INSERT INTO Playlists (name, user_id) VALUES (?, ?)",
		playlistName, user.Id)
	if err != nil {
		return c.Status(500).SendString("Помилка створення плейлиста")
	}

	// Отримання ID новоствореного плейлиста
	playlistID, err := result.LastInsertId()
	if err != nil {
		return c.Status(500).SendString("Помилка отримання ID плейлиста")
	}

	// Додавання продукту до плейлиста
	_, err = tx.Exec("INSERT INTO PlaylistItems (playlist_id, product_id) VALUES (?, ?)",
		playlistID, productID)
	if err != nil {
		return c.Status(500).SendString("Помилка додавання до плейлиста")
	}

	// Підтвердження транзакції
	if err = tx.Commit(); err != nil {
		return c.Status(500).SendString("Помилка підтвердження транзакції")
	}

	// Перенаправлення користувача до створеного плейлиста
	return c.Redirect("/playlist/" + fmt.Sprintf("%d", playlistID))
}
