package admin

import (
	"github.com/DjentBoiiii/marketplace/internal/auth"
	"github.com/DjentBoiiii/marketplace/internal/models"
	"github.com/DjentBoiiii/marketplace/internal/render"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gofiber/fiber/v2"
)

// ListProducts retrieves products filtered by type for admin management
// Displays all products of a specific type (audio, midi, samples) in the admin panel
func ListProducts(c *fiber.Ctx) error {
	user, err := auth.GetUserData(c)
	if err != nil {
		return c.Status(401).SendString("Необхідно увійти в систему")
	}

	productType := c.Query("type", "audio")

	rows, err := DB.Query(`
		SELECT p.id, p.name, p.price, p.type, p.description, p.vendor, p.genre, p.image_url, p.created_at
		FROM Products p
		WHERE p.type = ?
		ORDER BY p.created_at DESC`, productType)
	if err != nil {
		return c.Status(500).SendString("Помилка отримання списку продуктів")
	}
	defer rows.Close()

	var products []models.Product
	for rows.Next() {
		var product models.Product
		var createdAt []uint8
		err := rows.Scan(
			&product.ID,
			&product.Name,
			&product.Price,
			&product.Type,
			&product.Description,
			&product.Owner,
			&product.Genre,
			&product.ImageURL,
			&createdAt,
		)
		if err != nil {
			continue
		}
		products = append(products, product)
	}

	data := render.TemplateData{
		"products":    products,
		"user":        user,
		"productType": productType,
	}

	return render.RenderTemplate(c, "admin_products.html", data)
}
