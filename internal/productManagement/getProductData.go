package productManagement

import (
	"github.com/DjentBoiiii/marketplace/internal/models"
)

func GetProductData(productName, owner string) (models.Product, error) {
	var p models.Product
	var imagePath string
	err := DB.QueryRow(`
		SELECT p.name, p.price, p.type, u.name, p.product_img, p.description
		FROM Products p
		JOIN Users u ON p.vendor = u.id
		WHERE p.name = ? AND u.name = ?`, productName, owner).Scan(&p.Name, &p.Price, &p.Type, &p.Owner, &imagePath, &p.Description)
	if err != nil {
		return p, err
	}

	// Формуємо повний шлях до зображення
	p.ImageURL = imagePath
	return p, nil
}
