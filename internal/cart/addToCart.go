package cart

import (
	"github.com/DjentBoiiii/marketplace/internal/auth"
	"github.com/gofiber/fiber/v2"
)

func AddToCart(c *fiber.Ctx) error {
	user, err := auth.GetUserData(c)
	if err != nil {
		return c.Status(500).SendString("Помилка отримання даних користувача")
	}

	if user.Username == "" {
		return c.Redirect("/login")
	}

	productID := c.FormValue("product_id")
	if productID == "" {
		return c.Status(400).SendString("Не вказано ID продукту")
	}

	var count int
	err = DB.QueryRow("SELECT COUNT(*) FROM Cart WHERE user_id = ? AND product_id = ?", user.Id, productID).Scan(&count)
	if err != nil {
		return c.Status(500).SendString("Помилка перевірки кошика")
	}

	if count > 0 {
		return c.Status(400).SendString("Цей продукт вже у вашому кошику")
	}
	_, err = DB.Exec("INSERT INTO Cart (user_id, product_id) VALUES (?, ?)", user.Id, productID)
	if err != nil {
		return c.Status(500).SendString("Помилка додавання до кошика")
	}

	referer := c.Get("Referer")
	if referer == "" {
		referer = "/"
	}
	return c.Redirect(referer)
}
