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

		return render.RenderTemplate(c, "catalogue_filtered.html",
			[2]interface{}{"user", user},
			[2]interface{}{"productsByType", productsByType},
			[2]interface{}{"filters", filterOptions},
			[2]interface{}{"genres", genres},
			[2]interface{}{"vendors", vendors},
			[2]interface{}{"minPriceAll", minPriceAll},
			[2]interface{}{"maxPriceAll", maxPriceAll},
			[2]interface{}{"isFiltered", isFiltered},
			[2]interface{}{"search", search},
		)
	}

	products, err := productManagement.GetAllRandomProducts()
	if err != nil {
		fmt.Println("Error fetching products:", err)
		return c.Status(500).SendString("Помилка завантаження продуктів")
	}

	return render.RenderTemplate(c, "catalogue.html",
		[2]interface{}{"user", user},
		[2]interface{}{"audioProducts", products["audio"]},
		[2]interface{}{"midiProducts", products["midi"]},
		[2]interface{}{"samplesProducts", products["samples"]},
		[2]interface{}{"genres", genres},
		[2]interface{}{"vendors", vendors},
		[2]interface{}{"minPriceAll", minPriceAll},
		[2]interface{}{"maxPriceAll", maxPriceAll},
		[2]interface{}{"showAllTypesLink", true},
	)
}
