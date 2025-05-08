package productManagement

import (
	"fmt"
	"time"

	"github.com/DjentBoiiii/marketplace/internal/models"
)

func GetAllProductsByVendor(vendorUsername string) (map[string][]models.Product, models.VendorInfo, error) {

	// Get vendor info
	var vendorInfo models.VendorInfo
	vendorQuery := `SELECT username, profile_photo FROM Users WHERE username = ?`
	err := DB.QueryRow(vendorQuery, vendorUsername).Scan(&vendorInfo.Username, &vendorInfo.ProfilePic)
	if err != nil {
		return nil, models.VendorInfo{}, fmt.Errorf("failed to get vendor info: %v", err)
	}

	// Get all products by vendor
	query := `SELECT id, name, price, type, description, vendor, genre, image_url, created_at 
              FROM Products 
              WHERE vendor = ?
              ORDER BY type, created_at DESC`

	rows, err := DB.Query(query, vendorUsername)
	if err != nil {
		return nil, vendorInfo, fmt.Errorf("failed to query products: %v", err)
	}
	defer rows.Close()

	productsByType := make(map[string][]models.Product)
	for rows.Next() {
		var product models.Product
		var createdAt []uint8 // temporary variable for timestamp
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
		if err != nil {
			return nil, vendorInfo, fmt.Errorf("failed to scan product: %v", err)
		}

		// Parse the timestamp
		product.CreatedAt, err = time.Parse("2006-01-02 15:04:05", string(createdAt))
		if err != nil {
			return nil, vendorInfo, fmt.Errorf("failed to parse timestamp: %v", err)
		}

		// Append to corresponding type slice
		productsByType[product.Type] = append(productsByType[product.Type], product)
	}

	return productsByType, vendorInfo, nil
}
