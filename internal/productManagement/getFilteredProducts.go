package productManagement

import (
	"fmt"
	"time"

	"github.com/DjentBoiiii/marketplace/internal/models"
)

// FilterOptions defines the available product filter options
type FilterOptions struct {
	Name        string
	Vendor      string
	Genre       string
	MinPrice    int
	MaxPrice    int
	SortBy      string // Options: "name", "price_asc", "price_desc", "newest", "vendor"
	ProductType string
}

// GetFilteredProducts returns products filtered by the provided options
func GetFilteredProducts(options FilterOptions) ([]models.Product, error) {

	// Build the query
	queryBase := `SELECT id, name, price, type, description, vendor, genre, image_url, created_at 
                  FROM Products WHERE 1=1`

	var params []interface{}
	var conditions []string

	// Add type filter if provided
	if options.ProductType != "" {
		conditions = append(conditions, "type = ?")
		params = append(params, options.ProductType)
	}

	// Add name filter if provided
	if options.Name != "" {
		conditions = append(conditions, "name LIKE ?")
		params = append(params, "%"+options.Name+"%")
	}

	// Add vendor filter if provided
	if options.Vendor != "" {
		conditions = append(conditions, "vendor LIKE ?")
		params = append(params, "%"+options.Vendor+"%")
	}

	// Add genre filter if provided
	if options.Genre != "" {
		conditions = append(conditions, "genre LIKE ?")
		params = append(params, "%"+options.Genre+"%")
	}

	// Add price range filter if provided
	if options.MinPrice > 0 {
		conditions = append(conditions, "price >= ?")
		params = append(params, options.MinPrice)
	}
	if options.MaxPrice > 0 {
		conditions = append(conditions, "price <= ?")
		params = append(params, options.MaxPrice)
	}

	// Add conditions to query
	for _, condition := range conditions {
		queryBase += " AND " + condition
	}

	// Add sorting
	switch options.SortBy {
	case "name":
		queryBase += " ORDER BY name ASC"
	case "price_asc":
		queryBase += " ORDER BY price ASC"
	case "price_desc":
		queryBase += " ORDER BY price DESC"
	case "newest":
		queryBase += " ORDER BY created_at DESC"
	case "vendor":
		queryBase += " ORDER BY vendor ASC"
	default:
		queryBase += " ORDER BY created_at DESC" // Default sort by newest
	}

	// Execute the query
	stmt, err := DB.Prepare(queryBase)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare query: %v", err)
	}
	defer stmt.Close()

	rows, err := stmt.Query(params...)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %v", err)
	}
	defer rows.Close()

	// Parse the results
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

// GetGenres returns a list of all distinct genres in the database
func GetGenres() ([]string, error) {

	query := "SELECT DISTINCT genre FROM Products WHERE genre IS NOT NULL AND genre != ''"
	rows, err := DB.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query genres: %v", err)
	}
	defer rows.Close()

	var genres []string
	for rows.Next() {
		var genre string
		if err := rows.Scan(&genre); err != nil {
			return nil, fmt.Errorf("failed to scan genre: %v", err)
		}
		genres = append(genres, genre)
	}

	return genres, nil
}

// GetVendors returns a list of all distinct vendors in the database
func GetVendors() ([]string, error) {

	query := "SELECT DISTINCT vendor FROM Products"
	rows, err := DB.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query vendors: %v", err)
	}
	defer rows.Close()

	var vendors []string
	for rows.Next() {
		var vendor string
		if err := rows.Scan(&vendor); err != nil {
			return nil, fmt.Errorf("failed to scan vendor: %v", err)
		}
		vendors = append(vendors, vendor)
	}

	return vendors, nil
}

// GetPriceRange returns the minimum and maximum product prices in the database
func GetPriceRange(productType string) (int, int, error) {

	query := "SELECT MIN(price), MAX(price) FROM Products"
	params := []interface{}{}

	if productType != "" {
		query += " WHERE type = ?"
		params = append(params, productType)
	}

	var minPrice, maxPrice int
	err := DB.QueryRow(query, params...).Scan(&minPrice, &maxPrice)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to get price range: %v", err)
	}

	return minPrice, maxPrice, nil
}
