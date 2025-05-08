package cart

import (
	"fmt"

	"github.com/DjentBoiiii/marketplace/internal/auth"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gofiber/fiber/v2"
)

func PurchaseItems(c *fiber.Ctx) error {
	user, err := auth.GetUserData(c)
	if err != nil {
		return c.Status(500).SendString("Помилка отримання даних користувача")
	}

	cartProducts, err := GetCartProducts(user.Id)
	if err != nil {
		return c.Status(500).SendString("Помилка отримання продуктів кошика")
	}

	if len(cartProducts) == 0 {
		return c.Status(400).SendString("Кошик порожній")
	}

	tx, err := DB.Begin()
	if err != nil {
		return c.Status(500).SendString("Помилка створення транзакції")
	}

	for _, product := range cartProducts {

		var count int
		err = tx.QueryRow("SELECT COUNT(*) FROM Purchases WHERE user_id = ? AND product_id = ?", user.Id, product.ID).Scan(&count)
		if err != nil {
			tx.Rollback()
			return c.Status(500).SendString("Помилка перевірки покупок")
		}

		if count > 0 {
			continue
		}

		_, err = tx.Exec("INSERT INTO Purchases (user_id, product_id) VALUES (?, ?)", user.Id, product.ID)
		if err != nil {
			tx.Rollback()
			return c.Status(500).SendString(fmt.Sprintf("Помилка покупки продукту %s: %v", product.Name, err))
		}
	}

	_, err = tx.Exec("DELETE FROM Cart WHERE user_id = ?", user.Id)
	if err != nil {
		tx.Rollback()
		return c.Status(500).SendString("Помилка очищення кошика")
	}

	if err = tx.Commit(); err != nil {
		return c.Status(500).SendString("Помилка завершення транзакції")
	}

	return c.Redirect("/purchases")
}
