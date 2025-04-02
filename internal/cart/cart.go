package cart

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
func SetupCartHandlers(app *fiber.App) {
	app.Post("/cart/add", auth.LoginRequired(), AddToCart)
	app.Post("/cart/remove", auth.LoginRequired(), RemoveFromCart)
	app.Post("/cart/purchase", auth.LoginRequired(), PurchaseItems)
	app.Get("/cart", auth.LoginRequired(), ViewCart)
}
