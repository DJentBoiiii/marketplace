package comments

import (
	"github.com/DjentBoiiii/marketplace/internal/auth"
	"github.com/DjentBoiiii/marketplace/internal/db"
	"github.com/gofiber/fiber/v2"
)

var (
	DB = db.DB
)

func SetupCommentRoutes(app *fiber.App) {
	app.Post("/api/comments", auth.LoginRequired(), AddComment)
	app.Get("/api/comments/:id", GetComments)
	app.Delete("/api/comments/:id", auth.LoginRequired(), DeleteComment)
}
