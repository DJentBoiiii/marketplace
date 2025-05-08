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

		data := render.TemplateData{
			"product": product,
			"isOwned": isOwned,
			"user":    user,
		}

		return render.RenderTemplate(c, "product_info.html", data)
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

		data := render.TemplateData{
			"products": purchasedProducts,
			"user":     user,
		}

		return render.RenderTemplate(c, "purchases.html", data)
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

		data := render.TemplateData{
			"products":    products,
			"productType": productType,
			"user":        user,
		}

		return render.RenderTemplate(c, "purchases_category.html", data)
	})
}
