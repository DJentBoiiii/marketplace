// internal/playlist/playlist.go
package playlist

import (
	"database/sql"
	"os"

	"github.com/DjentBoiiii/marketplace/internal/auth"
	"github.com/DjentBoiiii/marketplace/internal/render"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

var (
	_           = godotenv.Load("/marketplace/.env")
	DB_USER     = os.Getenv("MYSQL_USER")
	DB_PASSWORD = os.Getenv("MYSQL_PASSWORD")
	DB_NAME     = os.Getenv("MYSQL_DATABASE")
	JWT_SECRET  = os.Getenv("JWT_SECRET")
	SHA_SECRET  = os.Getenv("SHA_SECRET")
	DB_HOST     = os.Getenv("DB_HOST")
	DB          *sql.DB
)

func init() {
	var err error
	DB, err = sql.Open("mysql", DB_USER+":"+DB_PASSWORD+"@tcp("+DB_HOST+":3306)/"+DB_NAME)
	if err != nil {
		panic(err)
	}
}

func SetupPlaylistHandlers(app *fiber.App) {
	app.Get("/playlists", auth.LoginRequired(), func(c *fiber.Ctx) error {
		user, err := auth.GetUserData(c)
		if err != nil {
			return c.Status(500).SendString("Помилка отримання даних користувача")
		}

		playlists, err := GetUserPlaylists(user.Id)
		if err != nil {
			return c.Status(500).SendString("Помилка отримання плейлистів користувача")
		}

		return render.RenderTemplate(c, "playlists.html", [2]interface{}{"playlists", playlists}, [2]interface{}{"user", user})
	})

	// API endpoint to get user playlists as JSON - Put this BEFORE :id route to avoid conflicts
	app.Get("/playlist/get-user-playlists", auth.LoginRequired(), func(c *fiber.Ctx) error {
		user, err := auth.GetUserData(c)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"error": "Помилка отримання даних користувача",
			})
		}

		playlists, err := GetUserPlaylists(user.Id)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"error": "Помилка отримання плейлистів користувача",
			})
		}

		return c.JSON(playlists)
	})

	// Add to playlist route
	app.Post("/playlist/add-item", auth.LoginRequired(), AddToPlaylist)

	// Create and add to playlist route
	app.Post("/playlist/create-and-add", auth.LoginRequired(), CreateAndAddToPlaylist)

	// Create playlist route
	app.Post("/playlist/create", auth.LoginRequired(), CreatePlaylist)

	// Delete playlist route
	app.Post("/playlist/delete", auth.LoginRequired(), DeletePlaylist)

	// Remove from playlist route
	app.Post("/playlist/remove", auth.LoginRequired(), RemoveFromPlaylist)

	app.Get("/playlist/:id", auth.LoginRequired(), func(c *fiber.Ctx) error {
		user, err := auth.GetUserData(c)
		if err != nil {
			return c.Status(500).SendString("Помилка отримання даних користувача")
		}

		playlistId := c.Params("id")

		// Перевірка власності плейлиста
		var userId int
		err = DB.QueryRow("SELECT user_id FROM Playlists WHERE id = ?", playlistId).Scan(&userId)
		if err != nil {
			return c.Status(404).SendString("Плейлист не знайдено")
		}

		if userId != user.Id {
			return c.Status(403).SendString("У вас немає доступу до цього плейлиста")
		}

		// Отримання даних плейлиста
		var playlist struct {
			Id    string
			Name  string
			Items []interface{}
		}
		playlist.Id = playlistId

		err = DB.QueryRow("SELECT name FROM Playlists WHERE id = ?", playlistId).Scan(&playlist.Name)
		if err != nil {
			return c.Status(500).SendString("Помилка отримання даних плейлиста")
		}

		// Отримання елементів плейлиста
		items, err := GetPlaylistItems(playlistId)
		if err != nil {
			return c.Status(500).SendString("Помилка отримання елементів плейлиста")
		}

		return render.RenderTemplate(c, "playlist_detail.html",
			[2]interface{}{"playlist", playlist},
			[2]interface{}{"items", items},
			[2]interface{}{"user", user})
	})
}
