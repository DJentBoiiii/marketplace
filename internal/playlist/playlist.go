// internal/playlist/playlist.go
package playlist

import (
	"database/sql"
	"os"

	"github.com/DjentBoiiii/marketplace/internal/auth"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

var (
	_           = godotenv.Load("/marketplace/.env")
	DB_USER     = os.Getenv("MYSQL_USER")
	DB_PASSWORD = os.Getenv("MYSQL_PASSWORD")
	DB_NAME     = os.Getenv("MYSQL_DATABASE")
	DB          *sql.DB
)

func init() {
	var err error
	DB, err = sql.Open("mysql", DB_USER+":"+DB_PASSWORD+"@tcp(boku-no-sukele:3306)/"+DB_NAME)
	if err != nil {
		panic(err)
	}
}

func SetupPlaylistHandlers(app *fiber.App) {
	app.Get("/playlists", auth.LoginRequired(), ViewPlaylists)
	app.Get("/playlist/:id", auth.LoginRequired(), ViewPlaylist)
	app.Post("/playlist/create", auth.LoginRequired(), CreatePlaylist)
	app.Post("/playlist/add", auth.LoginRequired(), AddToPlaylist)
	app.Post("/playlist/remove", auth.LoginRequired(), RemoveFromPlaylist)
	app.Post("/playlist/delete", auth.LoginRequired(), DeletePlaylist)
}
