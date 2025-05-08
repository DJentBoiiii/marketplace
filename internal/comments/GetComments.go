package comments

import (
	"github.com/DjentBoiiii/marketplace/internal/models"
	"github.com/gofiber/fiber/v2"
)

func GetComments(c *fiber.Ctx) error {
	productID := c.Params("id")
	if productID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "ID продукту не вказано",
		})
	}

	rows, err := DB.Query(`
		SELECT c.id, c.user_id, u.username, c.product_id, c.comment, c.likes_product, c.created_at, 
		       IFNULL(u.profile_photo, '') as profile_photo
		FROM Comments c
		JOIN Users u ON c.user_id = u.id
		WHERE c.product_id = ?
		ORDER BY c.created_at DESC`, productID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Помилка отримання коментарів",
		})
	}
	defer rows.Close()

	var comments []models.Comment
	for rows.Next() {
		var comment models.Comment
		err := rows.Scan(
			&comment.ID, &comment.UserID, &comment.Username, &comment.ProductID,
			&comment.Comment, &comment.LikesProduct, &comment.CreatedAt, &comment.ProfilePhoto,
		)
		if err != nil {
			continue
		}
		comments = append(comments, comment)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success":  true,
		"comments": comments,
	})
}
