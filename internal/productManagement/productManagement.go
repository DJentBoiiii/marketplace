package productManagement

import (
	"github.com/DjentBoiiii/marketplace/internal/auth"
	"github.com/DjentBoiiii/marketplace/internal/db"
	"github.com/DjentBoiiii/marketplace/internal/render"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gofiber/fiber/v2"
)

var (
	DB = db.DB
)

func SetupProductHandlers(app *fiber.App) {

	app.Get("/product/:name/:owner", func(c *fiber.Ctx) error {
		name := c.Params("name")
		owner := c.Params("owner")

		product, err := GetProductData(name, owner)
		if err != nil {
			return c.Status(404).SendString("Продукт не знайдено")
		}

		user, _ := auth.GetUserData(c)

		isOwned := false
		if user.Id > 0 {
			isOwned, _ = CheckUserOwnsProduct(user.Id, product.ID)
		}

		return render.RenderTemplate(c, "product_info.html",
			[2]interface{}{"product", product},
			[2]interface{}{"isOwned", isOwned},
			[2]interface{}{"user", user})
	})

	app.Get("/purchases", auth.LoginRequired(), func(c *fiber.Ctx) error {
		user, err := auth.GetUserData(c)
		if err != nil {
			return c.Status(500).SendString("Помилка отримання даних користувача")
		}

		purchasedProducts, err := ViewPurchases(user.Id)
		if err != nil {
			return c.Status(500).SendString("Помилка отримання придбаних продуктів")
		}

		return render.RenderTemplate(c, "purchases.html",
			[2]interface{}{"products", purchasedProducts},
			[2]interface{}{"user", user})
	})

	app.Get("/purchases/:type", auth.LoginRequired(), func(c *fiber.Ctx) error {
		productType := c.Params("type")
		if productType != "audio" && productType != "midi" && productType != "samples" {
			return c.Status(400).SendString("Невірний тип продукту")
		}

		user, err := auth.GetUserData(c)
		if err != nil {
			return c.Status(500).SendString("Помилка отримання даних користувача")
		}

		products, err := GetUserOwnedProducts(user.Id, productType)
		if err != nil {
			return c.Status(500).SendString("Помилка отримання придбаних продуктів")
		}

		return render.RenderTemplate(c, "purchases_category.html",
			[2]interface{}{"products", products},
			[2]interface{}{"productType", productType},
			[2]interface{}{"user", user})
	})
}
