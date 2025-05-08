package admin

import (
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gofiber/fiber/v2"
)

// DeleteProduct removes a product and all its associated data
// Deletes the product and cascades deletion to purchases, cart items, playlist items, and comments
func DeleteProduct(c *fiber.Ctx) error {
	productId := c.Params("id")
	productType := c.Query("type", "audio")

	var imagePath string
	_ = DB.QueryRow("SELECT image_url FROM Products WHERE id = ?", productId).Scan(&imagePath)

	_, err := DB.Exec("DELETE FROM Products WHERE id = ?", productId)
	if err != nil {
		return c.Status(500).SendString("Помилка видалення продукту")
	}

	_, err = DB.Exec("DELETE FROM Purchases WHERE product_id = ?", productId)
	if err != nil {

		fmt.Println("Error deleting purchases:", err)
	}

	_, err = DB.Exec("DELETE FROM Cart WHERE product_id = ?", productId)
	if err != nil {
		fmt.Println("Error deleting cart items:", err)
	}

	_, err = DB.Exec("DELETE FROM PlaylistItems WHERE product_id = ?", productId)
	if err != nil {
		fmt.Println("Error deleting playlist items:", err)
	}

	_, err = DB.Exec("DELETE FROM Comments WHERE product_id = ?", productId)
	if err != nil {
		fmt.Println("Error deleting comments:", err)
	}

	return c.Redirect("/admin/products?type=" + productType)
}
