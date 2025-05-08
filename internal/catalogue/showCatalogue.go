package catalogue

import (
	"fmt"
	"strconv"

	"github.com/DjentBoiiii/marketplace/internal/auth"
	"github.com/DjentBoiiii/marketplace/internal/productManagement"
	"github.com/DjentBoiiii/marketplace/internal/render"
	"github.com/gofiber/fiber/v2"
)

func ShowCatalogue(c *fiber.Ctx) error {
	user, _ := auth.GetUserData(c)

	// Get search query from the search form
	search := c.Query("search", "")

	name := c.Query("name", "")
	vendor := c.Query("vendor", "")
	genre := c.Query("genre", "")
	minPriceStr := c.Query("min_price", "")
	maxPriceStr := c.Query("max_price", "")
	sortBy := c.Query("sort", "newest")

	// If there's a search query, use it as the name filter
	if search != "" {
		name = search
	}

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
	minPriceAll, maxPriceAll, _ := productManagement.GetPriceRange("")

	if isFiltered {

		filterOptions := productManagement.FilterOptions{
			Name:     name,
			Vendor:   vendor,
			Genre:    genre,
			MinPrice: minPrice,
			MaxPrice: maxPrice,
			SortBy:   sortBy,
		}

		filteredProducts, err := productManagement.GetFilteredProducts(filterOptions)
		if err != nil {
			fmt.Println("Error fetching filtered products:", err)
			return c.Status(500).SendString("Помилка завантаження продуктів")
		}

		productsByType := make(map[string][]interface{})
		for _, product := range filteredProducts {
			productsByType[product.Type] = append(productsByType[product.Type], product)
		}

		data := render.TemplateData{
			"user":           user,
			"productsByType": productsByType,
			"filters":        filterOptions,
			"genres":         genres,
			"vendors":        vendors,
			"minPriceAll":    minPriceAll,
			"maxPriceAll":    maxPriceAll,
			"isFiltered":     isFiltered,
			"search":         search,
		}

		return render.RenderTemplate(c, "catalogue_filtered.html", data)
	}

	products, err := productManagement.GetAllRandomProducts()
	if err != nil {
		fmt.Println("Error fetching products:", err)
		return c.Status(500).SendString("Помилка завантаження продуктів")
	}

	data := render.TemplateData{
		"user":             user,
		"audioProducts":    products["audio"],
		"midiProducts":     products["midi"],
		"samplesProducts":  products["samples"],
		"genres":           genres,
		"vendors":          vendors,
		"minPriceAll":      minPriceAll,
		"maxPriceAll":      maxPriceAll,
		"showAllTypesLink": true,
	}

	return render.RenderTemplate(c, "catalogue.html", data)
}
