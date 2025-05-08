package comments

import (
	"fmt"
	"time"

	"github.com/DjentBoiiii/marketplace/internal/auth"
	"github.com/DjentBoiiii/marketplace/internal/models"
	"github.com/gofiber/fiber/v2"
)

func AddComment(c *fiber.Ctx) error {
	fmt.Println("AddComment function called") // Debug log

	user, err := auth.GetUserData(c)
	if err != nil {
		fmt.Println("Error getting user data:", err) // Debug log
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": "Помилка отримання даних користувача",
		})
	}
	fmt.Println("User data retrieved:", user.Id, user.Username) // Debug log

	var input struct {
		ProductID    int    `json:"product_id"`
		Comment      string `json:"comment"`
		LikesProduct bool   `json:"likes_product"`
	}

	if err := c.BodyParser(&input); err != nil {
		fmt.Println("Error parsing request body:", err) // Debug log
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"success": false,
			"message": "Невірний формат даних запиту",
		})
	}
	fmt.Println("Comment data parsed:", input.ProductID, input.Comment, input.LikesProduct) // Debug log

	if input.ProductID <= 0 || input.Comment == "" {
		fmt.Println("Invalid comment data") // Debug log
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Невірні дані коментаря",
		})
	}

	// Debug: Print the SQL query being executed
	fmt.Println("Executing INSERT query with values:", user.Id, input.ProductID, input.Comment, input.LikesProduct)

	result, err := DB.Exec(
		"INSERT INTO Comments (user_id, product_id, comment, likes_product) VALUES (?, ?, ?, ?)",
		user.Id, input.ProductID, input.Comment, input.LikesProduct,
	)
	if err != nil {
		fmt.Println("Error executing insert query:", err) // Debug log
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Помилка додавання коментаря",
		})
	}

	commentID, err := result.LastInsertId()
	if err != nil {
		fmt.Println("Error getting last insert ID:", err) // Debug log
	}
	fmt.Println("Comment inserted with ID:", commentID) // Debug log

	var comment models.Comment
	var createdAtStr string // Use string as intermediate type for date

	err = DB.QueryRow(`
		SELECT c.id, c.user_id, u.username, c.product_id, c.comment, c.likes_product, c.created_at, u.profile_photo 
		FROM Comments c
		JOIN Users u ON c.user_id = u.id
		WHERE c.id = ?`, commentID).Scan(
		&comment.ID, &comment.UserID, &comment.Username, &comment.ProductID,
		&comment.Comment, &comment.LikesProduct, &createdAtStr, &comment.ProfilePhoto,
	)

	if err != nil {
		fmt.Println("Error retrieving inserted comment:", err) // Debug log
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"success":    true,
			"message":    "Коментар додано",
			"comment_id": commentID,
		})
	}

	// Parse the date string into time.Time
	parsedTime, err := time.Parse("2006-01-02 15:04:05", createdAtStr)
	if err != nil {
		fmt.Println("Error parsing date string:", err, "Date string:", createdAtStr)
		// Use current time as fallback
		comment.CreatedAt = time.Now()
	} else {
		comment.CreatedAt = parsedTime
	}

	fmt.Println("Comment successfully added and retrieved") // Debug log
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Коментар додано",
		"comment": comment,
	})
}
