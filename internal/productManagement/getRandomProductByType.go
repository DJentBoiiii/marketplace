package productManagement

import (
	"fmt"

	"time"

	"github.com/DjentBoiiii/marketplace/internal/models"
	_ "github.com/go-sql-driver/mysql"
)

func GetRandomProductsByType(productType string) ([]models.Product, error) {

	// Use ORDER BY RAND() to get random products
	query := `SELECT id, name, price, type, description, vendor, genre, image_url, created_at 
              FROM Products 
              WHERE type = ? 
              ORDER BY RAND() 
              LIMIT 10`

	rows, err := DB.Query(query, productType)
	if err != nil {
		return nil, fmt.Errorf("failed to query products: %v", err)
	}
	defer rows.Close()

	products := []models.Product{}
	for rows.Next() {
		var product models.Product
		var createdAt string
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
		if createdAt != "" {
			product.CreatedAt, err = time.Parse("2006-01-02 15:04:05", createdAt)
			if err != nil {
				return nil, fmt.Errorf("failed to parse created_at time: %v", err)
			}
		}
		if err != nil {
			return nil, fmt.Errorf("failed to scan product: %v", err)
		}
		products = append(products, product)
	}

	return products, nil
}
