package likes

import (
	"github.com/DjentBoiiii/marketplace/internal/auth"
	"github.com/DjentBoiiii/marketplace/internal/db"
	"github.com/gofiber/fiber/v2"
)

var (
	DB = db.DB
)

func SetupLikeRoutes(app *fiber.App) {
	app.Post("/api/likes/:id", auth.LoginRequired(), LikeProduct)
	app.Delete("/api/likes/:id", auth.LoginRequired(), UnlikeProduct)
	app.Get("/api/likes/:id", GetLikeStatus)
	app.Get("/api/likes/:id/count", GetLikeCount)
}
