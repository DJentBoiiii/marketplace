package catalogue

import (
	"fmt"
	"strconv"

	"github.com/DjentBoiiii/marketplace/internal/auth"
	"github.com/DjentBoiiii/marketplace/internal/productManagement"
	"github.com/DjentBoiiii/marketplace/internal/render"
	"github.com/gofiber/fiber/v2"
)

func ShowProductsByType(c *fiber.Ctx) error {
	productType := c.Params("type")
	user, _ := auth.GetUserData(c)

	validTypes := map[string]bool{
		"audio":   true,
		"midi":    true,
		"samples": true,
	}

	if !validTypes[productType] {
		return c.Status(404).SendString("Невідомий тип продукту")
	}

	name := c.Query("name", "")
	vendor := c.Query("vendor", "")
	genre := c.Query("genre", "")
	minPriceStr := c.Query("min_price", "")
	maxPriceStr := c.Query("max_price", "")
	sortBy := c.Query("sort", "newest")

	var minPrice, maxPrice int
	if minPriceStr != "" {
		minPrice, _ = strconv.Atoi(minPriceStr)
	}
	if maxPriceStr != "" {
		maxPrice, _ = strconv.Atoi(maxPriceStr)
	}

	isFiltered := name != "" || vendor != "" || genre != "" || minPrice > 0 || maxPrice > 0 || sortBy != "newest"

	genres, _ := productManagement.GetGenres()
	vendors, _ := productManagement.GetVendors()
	minPriceRange, maxPriceRange, _ := productManagement.GetPriceRange(productType)

	filterOptions := productManagement.FilterOptions{
		Name:        name,
		Vendor:      vendor,
		Genre:       genre,
		MinPrice:    minPrice,
		MaxPrice:    maxPrice,
		SortBy:      sortBy,
		ProductType: productType,
	}

	products, err := productManagement.GetFilteredProducts(filterOptions)
	if err != nil {
		fmt.Println("Error fetching products:", err)
		return c.Status(500).SendString("Помилка завантаження продуктів")
	}

	data := render.TemplateData{
		"products":      products,
		"user":          user,
		"type":          productType,
		"title":         formatTypeTitle(productType),
		"filters":       filterOptions,
		"genres":        genres,
		"vendors":       vendors,
		"minPriceRange": minPriceRange,
		"maxPriceRange": maxPriceRange,
		"isFiltered":    isFiltered,
	}

	return render.RenderTemplate(c, "catalogue_type.html", data)
}
