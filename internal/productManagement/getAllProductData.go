package productManagement

import (
	"github.com/DjentBoiiii/marketplace/internal/models"
)

func GetAllProductsData(owner, productType string) ([]models.Product, error) {
	rows, err := DB.Query(`
		SELECT p.name, p.price, p.type, u.username, p.image_url 
		FROM Products p
		JOIN Users u ON p.vendor = u.username
		WHERE u.username = ? AND p.type = ?`, owner, productType)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []models.Product
	for rows.Next() {
		var p models.Product
		var imagePath string
		if err := rows.Scan(&p.Name, &p.Price, &p.Type, &p.Owner, &imagePath); err != nil {
			return nil, err
		}
		p.ImageURL = "/" + imagePath
		products = append(products, p)
	}
	return products, nil
}
