package likes

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
)

// GetLikeCount returns the number of likes for a specific product
func GetLikeCount(c *fiber.Ctx) error {
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

	// Check if the product exists
	var exists bool
	err = DB.QueryRow("SELECT EXISTS (SELECT 1 FROM Products WHERE id = ?)", productID).Scan(&exists)
	if err != nil || !exists {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": "Product not found",
		})
	}

	// Get total likes count
	var count int
	err = DB.QueryRow("SELECT COUNT(*) FROM ProductLikes WHERE product_id = ?", productID).Scan(&count)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Error getting likes count",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"count":   count,
	})
}
