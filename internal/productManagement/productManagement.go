package productManagement

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/DjentBoiiii/marketplace/internal/auth"
	"github.com/DjentBoiiii/marketplace/internal/playlist"
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
	DB          *sql.DB
)

func init() {
	var err error
	DB, err = sql.Open("mysql", DB_USER+":"+DB_PASSWORD+"@tcp(boku-no-sukele:3306)/"+DB_NAME)
	if err != nil {
		panic(err)
	}
}

func SetupProductHandlers(app *fiber.App) {
	// Каталог продуктів (всі продукти певного користувача і певного типу)
	app.Get("/catalogue/:owner/:type", func(c *fiber.Ctx) error {
		owner := c.Params("owner")
		productType := c.Params("type")

		products, err := GetAllProductsData(owner, productType)
		if err != nil {
			fmt.Println(err)
			return c.Status(500).SendString("Помилка отримання даних")
		}
		user, _ := auth.GetUserData(c)
		playlist, _ := playlist.GetUserPlaylists(user.Id)
		return render.RenderTemplate(c, "catalogue.html", [2]interface{}{"products", products}, [2]interface{}{"playlist", playlist})
	})

	// Інформація про один продукт
	app.Get("/product/:name/:owner", func(c *fiber.Ctx) error {
		name := c.Params("name")
		owner := c.Params("owner")

		product, err := GetProductData(name, owner)
		if err != nil {
			return c.Status(404).SendString("Продукт не знайдено")
		}

		return render.RenderTemplate(c, "product_info.html", [2]interface{}{"product", product})
	})
}
