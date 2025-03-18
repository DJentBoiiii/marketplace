package cart

import (
	"github.com/DjentBoiiii/marketplace/internal/auth"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gofiber/fiber/v2"
)

func RemoveFromCart(c *fiber.Ctx) error {
	user, err := auth.GetUserData(c)
	if err != nil {
		return c.Status(500).SendString("Помилка отримання даних користувача")
	}

	productID := c.FormValue("product_id")
	if productID == "" {
		return c.Status(400).SendString("Не вказано ID продукту")
	}

	_, err = DB.Exec("DELETE FROM Cart WHERE user_id = ? AND product_id = ?", user.Id, productID)
	if err != nil {
		return c.Status(500).SendString("Помилка видалення з кошика")
	}

	return c.Redirect("/cart")
}
