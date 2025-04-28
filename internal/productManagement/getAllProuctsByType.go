package productManagement

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/DjentBoiiii/marketplace/internal/models"
	_ "github.com/go-sql-driver/mysql"
)

func GetAllProductsByType(productType string) ([]models.Product, error) {
	db, err := sql.Open("mysql", DB_USER+":"+DB_PASSWORD+"@tcp("+DB_HOST+":3306)/"+DB_NAME)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %v", err)
	}
	defer db.Close()

	query := `SELECT id, name, price, type, description, vendor, genre, image_url, created_at 
              FROM Products 
              WHERE type = ? 
              ORDER BY created_at DESC`

	rows, err := db.Query(query, productType)
	if err != nil {
		return nil, fmt.Errorf("failed to query products: %v", err)
	}
	defer rows.Close()

	products := []models.Product{}
	for rows.Next() {
		var product models.Product
		var createdAt string
		err := rows.Scan(
			&product.Id,
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
