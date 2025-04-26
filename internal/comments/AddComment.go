package comments

import (
	"database/sql"
	"fmt"

	"github.com/DjentBoiiii/marketplace/internal/auth"
	"github.com/DjentBoiiii/marketplace/internal/models"
	"github.com/gofiber/fiber/v2"
)

func AddComment(c *fiber.Ctx) error {
	user, err := auth.GetUserData(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": "Помилка отримання даних користувача",
		})
	}

	var input struct {
		ProductID    int    `json:"product_id"`
		Comment      string `json:"comment"`
		LikesProduct bool   `json:"likes_product"`
	}

	if err := c.BodyParser(&input); err != nil {
		fmt.Println(err)
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"success": false,
			"message": "Невірний формат даних запиту",
		})
	}

	if input.ProductID <= 0 || input.Comment == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Невірні дані коментаря",
		})
	}

	db, err := sql.Open("mysql", DB_USER+":"+DB_PASSWORD+"@tcp("+DB_HOST+":3306)/"+DB_NAME)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Помилка підключення до бази даних",
		})
	}
	defer db.Close()

	result, err := db.Exec(
		"INSERT INTO Comments (user_id, product_id, comment, likes_product) VALUES (?, ?, ?, ?)",
		user.Id, input.ProductID, input.Comment, input.LikesProduct,
	)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Помилка додавання коментаря",
		})
	}

	commentID, _ := result.LastInsertId()

	var comment models.Comment
	err = db.QueryRow(`
		SELECT c.id, c.user_id, u.username, c.product_id, c.comment, c.likes_product, c.created_at, u.profile_photo 
		FROM Comments c
		JOIN Users u ON c.user_id = u.id
		WHERE c.id = ?`, commentID).Scan(
		&comment.ID, &comment.UserID, &comment.Username, &comment.ProductID,
		&comment.Comment, &comment.LikesProduct, &comment.CreatedAt, &comment.ProfilePhoto,
	)
	if err != nil {

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"success":    true,
			"message":    "Коментар додано",
			"comment_id": commentID,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Коментар додано",
		"comment": comment,
	})
}
