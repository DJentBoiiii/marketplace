package likes

import (
	"strconv"

	"github.com/DjentBoiiii/marketplace/internal/auth"
	"github.com/DjentBoiiii/marketplace/internal/models"
	"github.com/gofiber/fiber/v2"
)

// GetLikeStatus checks if a user has liked a specific product
func GetLikeStatus(c *fiber.Ctx) error {
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

	// Get current user data (if logged in)
	user, err := auth.GetUserData(c)

	// Default to not liked if user is not logged in
	isLiked := false

	// If user is logged in, check if they liked the product
	if err == nil {
		err = DB.QueryRow("SELECT EXISTS (SELECT 1 FROM ProductLikes WHERE user_id = ? AND product_id = ?)",
			user.Id, productID).Scan(&isLiked)

		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"success": false,
				"message": "Error checking like status",
			})
		}
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

	return c.JSON(models.ProductLikeResponse{
		Success: true,
		IsLiked: isLiked,
		Count:   count,
	})
}
