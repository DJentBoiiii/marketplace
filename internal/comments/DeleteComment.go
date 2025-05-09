package comments

import (
	"github.com/DjentBoiiii/marketplace/internal/auth"
	"github.com/gofiber/fiber/v2"
)

func DeleteComment(c *fiber.Ctx) error {
	user, err := auth.GetUserData(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": "Помилка отримання даних користувача",
		})
	}

	commentID := c.Params("id")
	if commentID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "ID коментаря не вказано",
		})
	}

	var commentUserID int
	var isAdmin bool
	err = DB.QueryRow("SELECT user_id FROM Comments WHERE id = ?", commentID).Scan(&commentUserID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": "Коментар не знайдено",
		})
	}

	err = DB.QueryRow("SELECT is_admin FROM Users WHERE id = ?", user.Id).Scan(&isAdmin)
	if err != nil {
		isAdmin = false
	}

	if commentUserID != user.Id && !isAdmin {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"success": false,
			"message": "Немає прав на видалення цього коментаря",
		})
	}

	_, err = DB.Exec("DELETE FROM Comments WHERE id = ?", commentID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Помилка видалення коментаря",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Коментар видалено",
	})
}
