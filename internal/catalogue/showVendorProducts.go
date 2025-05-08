package catalogue

import (
	"fmt"
	"strconv"

	"github.com/DjentBoiiii/marketplace/internal/auth"
	"github.com/DjentBoiiii/marketplace/internal/models"
	"github.com/DjentBoiiii/marketplace/internal/productManagement"
	"github.com/DjentBoiiii/marketplace/internal/render"
	"github.com/gofiber/fiber/v2"
)

func ShowVendorProducts(c *fiber.Ctx) error {
	vendorUsername := c.Params("username")
	user, _ := auth.GetUserData(c)

	productsByType, vendorInfo, err := productManagement.GetAllProductsByVendor(vendorUsername)
	if err != nil {
		fmt.Println("Error fetching vendor products:", err)
		return c.Status(500).SendString("Помилка завантаження продуктів")
	}

	name := c.Query("name", "")
	genre := c.Query("genre", "")
	minPriceStr := c.Query("min_price", "")
	maxPriceStr := c.Query("max_price", "")
	sortBy := c.Query("sort", "newest")
	productType := c.Query("type", "")

	var minPrice, maxPrice int
	if minPriceStr != "" {
		minPrice, _ = strconv.Atoi(minPriceStr)
	}
	if maxPriceStr != "" {
		maxPrice, _ = strconv.Atoi(maxPriceStr)
	}

	isFiltered := name != "" || genre != "" || minPrice > 0 || maxPrice > 0 || sortBy != "newest" || productType != ""

	if isFiltered {

		filterOptions := productManagement.FilterOptions{
			Name:        name,
			Vendor:      vendorUsername,
			Genre:       genre,
			MinPrice:    minPrice,
			MaxPrice:    maxPrice,
			SortBy:      sortBy,
			ProductType: productType,
		}

		filteredProducts, err := productManagement.GetFilteredProducts(filterOptions)
		if err != nil {
			fmt.Println("Error fetching filtered products:", err)
			return c.Status(500).SendString("Помилка завантаження продуктів")
		}

		filteredByType := make(map[string][]models.Product)
		for _, product := range filteredProducts {
			filteredByType[product.Type] = append(filteredByType[product.Type], product)
		}

		genres, _ := productManagement.GetGenres()
		minPriceRange, maxPriceRange, _ := productManagement.GetPriceRange("")

		data := render.TemplateData{
			"user":           user,
			"productsByType": filteredByType,
			"vendor":         vendorInfo,
			"filters":        filterOptions,
			"genres":         genres,
			"minPriceRange":  minPriceRange,
			"maxPriceRange":  maxPriceRange,
			"isFiltered":     isFiltered,
		}

		return render.RenderTemplate(c, "vendor_products.html", data)
	}

	genres, _ := productManagement.GetGenres()
	minPriceRange, maxPriceRange, _ := productManagement.GetPriceRange("")

	data := render.TemplateData{
		"user":           user,
		"productsByType": productsByType,
		"vendor":         vendorInfo,
		"genres":         genres,
		"minPriceRange":  minPriceRange,
		"maxPriceRange":  maxPriceRange,
	}

	return render.RenderTemplate(c, "vendor_products.html", data)
}
