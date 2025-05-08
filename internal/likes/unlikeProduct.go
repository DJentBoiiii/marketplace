package likes

import (
	"strconv"

	"github.com/DjentBoiiii/marketplace/internal/auth"
	"github.com/DjentBoiiii/marketplace/internal/models"
	"github.com/gofiber/fiber/v2"
)

// UnlikeProduct handles removing a like from a product
func UnlikeProduct(c *fiber.Ctx) error {
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

	// Check if user liked the product
	var liked bool
	err = DB.QueryRow("SELECT EXISTS (SELECT 1 FROM ProductLikes WHERE user_id = ? AND product_id = ?)",
		user.Id, productID).Scan(&liked)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Error checking like status",
		})
	}

	if !liked {
		// User had not liked this product
		var count int
		DB.QueryRow("SELECT COUNT(*) FROM ProductLikes WHERE product_id = ?", productID).Scan(&count)

		return c.JSON(models.ProductLikeResponse{
			Success: true,
			Message: "Product was not liked",
			IsLiked: false,
			Count:   count,
		})
	}

	// Remove the like
	_, err = DB.Exec(
		"DELETE FROM ProductLikes WHERE user_id = ? AND product_id = ?",
		user.Id, productID,
	)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Error removing like",
		})
	}

	// Get the updated count of likes
	var count int
	DB.QueryRow("SELECT COUNT(*) FROM ProductLikes WHERE product_id = ?", productID).Scan(&count)

	return c.JSON(models.ProductLikeResponse{
		Success: true,
		Message: "Product unliked successfully",
		IsLiked: false,
		Count:   count,
	})
}
