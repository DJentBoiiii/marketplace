package productManagement

import (
	"database/sql"

	"github.com/DjentBoiiii/marketplace/internal/models"
)

func CheckUserOwnsProduct(userId int, productId int) (bool, error) {
	var count int
	err := DB.QueryRow("SELECT COUNT(*) FROM Purchases WHERE user_id = ? AND product_id = ?", userId, productId).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func GetUserOwnedProducts(userId int, productType string) ([]models.Product, error) {
	query := `
		SELECT p.id, p.name, p.price, p.type, u.username, p.image_url, p.description, p.genre
		FROM Products p
		JOIN Users u ON p.vendor = u.username
		JOIN Purchases pur ON p.id = pur.product_id
		WHERE pur.user_id = ? AND p.type = ?`

	rows, err := DB.Query(query, userId, productType)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []models.Product
	for rows.Next() {
		var p models.Product
		var imagePath string
		if err := rows.Scan(&p.Id, &p.Name, &p.Price, &p.Type, &p.Owner, &imagePath, &p.Description, &p.Genre); err != nil {
			return nil, err
		}
		p.ImageURL = "/" + imagePath
		products = append(products, p)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return products, nil
}

func ViewPurchases(userId int) (map[string][]models.Product, error) {
	result := make(map[string][]models.Product)

	audioProducts, err := GetUserOwnedProducts(userId, "audio")
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	result["audio"] = audioProducts

	midiProducts, err := GetUserOwnedProducts(userId, "midi")
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	result["midi"] = midiProducts

	samplesProducts, err := GetUserOwnedProducts(userId, "samples")
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	result["samples"] = samplesProducts

	return result, nil
}
