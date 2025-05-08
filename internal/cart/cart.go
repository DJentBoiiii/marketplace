package cart

import (
	"github.com/DjentBoiiii/marketplace/internal/auth"
	"github.com/DjentBoiiii/marketplace/internal/db"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gofiber/fiber/v2"
)

var (
	DB = db.DB
)

func SetupCartHandlers(app *fiber.App) {
	app.Post("/cart/add", auth.LoginRequired(), AddToCart)
	app.Post("/cart/remove", auth.LoginRequired(), RemoveFromCart)
	app.Post("/cart/purchase", auth.LoginRequired(), PurchaseItems)
	app.Get("/cart", auth.LoginRequired(), ViewCart)
}
