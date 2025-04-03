package comments

import (
	"os"

	"github.com/DjentBoiiii/marketplace/internal/auth"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

var (
	_           = godotenv.Load("/marketplace/.env")
	DB_USER     = os.Getenv("MYSQL_USER")
	DB_PASSWORD = os.Getenv("MYSQL_PASSWORD")
	DB_NAME     = os.Getenv("MYSQL_DATABASE")
	DB_HOST     = os.Getenv("DB_HOST")
)

func SetupCommentRoutes(app *fiber.App) {
	app.Post("/api/comments", auth.LoginRequired(), AddComment)
	app.Get("/api/comments/:id", GetComments)
	app.Delete("/api/comments/:id", auth.LoginRequired(), DeleteComment)
}
