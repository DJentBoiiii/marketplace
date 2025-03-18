package cart

import (
	"github.com/DjentBoiiii/marketplace/internal/auth"
	"github.com/DjentBoiiii/marketplace/internal/render"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gofiber/fiber/v2"
)

func ViewCart(c *fiber.Ctx) error {
	user, err := auth.GetUserData(c)
	if err != nil {
		return c.Status(500).SendString("Помилка отримання даних користувача")
	}

	if user.Username == "" {
		return c.Redirect("/login")
	}

	products, err := GetCartProducts(user.Id)
	if err != nil {
		return c.Status(500).SendString("Помилка отримання продуктів кошика")
	}

	return render.RenderTemplate(c, "cart.html",
		[2]interface{}{"products", products},
		[2]interface{}{"user", user},
	)
}
