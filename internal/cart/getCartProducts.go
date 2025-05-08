package cart

import (
	"github.com/DjentBoiiii/marketplace/internal/models"
	_ "github.com/go-sql-driver/mysql"
)

func GetCartProducts(userId int) ([]models.Product, error) {
	rows, err := DB.Query(`
		SELECT p.id, p.name, p.price, p.type, u.username, p.image_url, p.description, p.genre
		FROM Products p
		JOIN Users u ON p.vendor = u.username
		JOIN Cart c ON p.id = c.product_id
		WHERE c.user_id = ?`, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []models.Product
	for rows.Next() {
		var p models.Product
		var imagePath string
		if err := rows.Scan(&p.ID, &p.Name, &p.Price, &p.Type, &p.Owner, &imagePath, &p.Description, &p.Genre); err != nil {
			return nil, err
		}
		p.ImageURL = "/" + imagePath
		products = append(products, p)
	}
	return products, nil
}
