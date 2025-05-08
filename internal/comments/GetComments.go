package comments

import (
	"fmt"
	"time"

	"github.com/DjentBoiiii/marketplace/internal/models"
	"github.com/gofiber/fiber/v2"
)

func GetComments(c *fiber.Ctx) error {
	productID := c.Params("id")
	fmt.Println("=== GetComments DEBUG ===")
	fmt.Println("Called for product ID:", productID)

	if productID == "" {
		fmt.Println("ERROR: Product ID not provided")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "ID продукту не вказано",
		})
	}

	// Print the full request details
	fmt.Println("Request URL:", c.OriginalURL())
	fmt.Println("Request method:", c.Method())
	fmt.Println("Request headers:", c.GetReqHeaders())

	// Debug: Print the SQL query being executed
	query := `
		SELECT c.id, c.user_id, u.username, c.product_id, c.comment, c.likes_product, c.created_at, 
		       IFNULL(u.profile_photo, '') as profile_photo
		FROM Comments c
		JOIN Users u ON c.user_id = u.id
		WHERE c.product_id = ?
		ORDER BY c.created_at DESC`

	fmt.Println("SQL query:", query)
	fmt.Println("With parameter:", productID)

	rows, err := DB.Query(query, productID)
	if err != nil {
		fmt.Println("DATABASE ERROR:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Помилка отримання коментарів",
			"error":   err.Error(),
		})
	}
	defer rows.Close()

	var comments []models.Comment
	rowCount := 0
	for rows.Next() {
		rowCount++
		var comment models.Comment
		var createdAtStr string // Use string as intermediate type for date

		err := rows.Scan(
			&comment.ID,
			&comment.UserID,
			&comment.Username,
			&comment.ProductID,
			&comment.Comment,
			&comment.LikesProduct,
			&createdAtStr, // Scan into string first
			&comment.ProfilePhoto,
		)

		if err != nil {
			fmt.Println("Error scanning row:", err)
			continue
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

		comments = append(comments, comment)
		fmt.Printf("Row %d: ID=%d, UserID=%d, Username=%s, Comment=%s, CreatedAt=%s\n",
			rowCount, comment.ID, comment.UserID, comment.Username, comment.Comment, comment.CreatedAt.String())
	}

	if err := rows.Err(); err != nil {
		fmt.Println("Error after row iteration:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Помилка читання коментарів",
			"error":   err.Error(),
		})
	}

	// Log the response we're sending back
	fmt.Println("Returning", len(comments), "comments")
	fmt.Println("=== END GetComments DEBUG ===")

	response := fiber.Map{
		"success":  true,
		"comments": comments,
	}

	// Debug the actual JSON being returned
	fmt.Println("Response:", response)

	return c.Status(fiber.StatusOK).JSON(response)
}
