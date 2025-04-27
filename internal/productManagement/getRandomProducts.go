package productManagement

import (
	"fmt"

	"github.com/DjentBoiiii/marketplace/internal/models"
	_ "github.com/go-sql-driver/mysql"
)

func GetAllRandomProducts() (map[string][]models.Product, error) {
	productTypes := []string{"audio", "midi", "samples"}
	result := make(map[string][]models.Product)

	for _, productType := range productTypes {
		products, err := GetRandomProductsByType(productType)
		if err != nil {
			return nil, fmt.Errorf("failed to get %s products: %v", productType, err)
		}
		result[productType] = products
	}

	return result, nil
}
