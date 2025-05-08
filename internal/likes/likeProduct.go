package likes

import (
	"strconv"

	"github.com/DjentBoiiii/marketplace/internal/auth"
	"github.com/DjentBoiiii/marketplace/internal/models"
	"github.com/gofiber/fiber/v2"
)

// LikeProduct handles adding a like to a product
func LikeProduct(c *fiber.Ctx) error {
	// Get product ID from URL
	productIDStr := c.Params("id")
	if productIDStr == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Product ID is required",
		})
	}

	productID, err := strconv.Atoi(productIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Invalid product ID",
		})
	}

	// Get current user data
	user, err := auth.GetUserData(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": "User authentication error",
		})
	}

	// Check if the product exists
	var exists bool
	err = DB.QueryRow("SELECT EXISTS (SELECT 1 FROM Products WHERE id = ?)", productID).Scan(&exists)
	if err != nil || !exists {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": "Product not found",
		})
	}

	// Check if user already liked the product
	var liked bool
	err = DB.QueryRow("SELECT EXISTS (SELECT 1 FROM ProductLikes WHERE user_id = ? AND product_id = ?)",
		user.Id, productID).Scan(&liked)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Error checking like status",
		})
	}

	if liked {
		// User already liked this product
		var count int
		DB.QueryRow("SELECT COUNT(*) FROM ProductLikes WHERE product_id = ?", productID).Scan(&count)

		return c.JSON(models.ProductLikeResponse{
			Success: true,
			Message: "Product already liked",
			IsLiked: true,
			Count:   count,
		})
	}

	// Add the like
	_, err = DB.Exec(
		"INSERT INTO ProductLikes (user_id, product_id) VALUES (?, ?)",
		user.Id, productID,
	)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Error adding like",
		})
	}

	// Get the updated count of likes
	var count int
	DB.QueryRow("SELECT COUNT(*) FROM ProductLikes WHERE product_id = ?", productID).Scan(&count)

	return c.JSON(models.ProductLikeResponse{
		Success: true,
		Message: "Product liked successfully",
		IsLiked: true,
		Count:   count,
	})
}
