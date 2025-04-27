package catalogue

import (
	"fmt"

	"github.com/DjentBoiiii/marketplace/internal/auth"
	"github.com/DjentBoiiii/marketplace/internal/productManagement"
	"github.com/DjentBoiiii/marketplace/internal/render"
	"github.com/gofiber/fiber/v2"
)

// SetupCatalogueRoutes configures the catalogue related routes
func SetupCatalogueRoutes(app *fiber.App) {
	app.Get("/catalogue", ShowCatalogue)
	app.Get("/catalogue/:type", ShowProductsByType)
	app.Get("/test/vendor/:username", ShowVendorProducts)
}

// ShowCatalogue displays the first 10 random products of each type
func ShowCatalogue(c *fiber.Ctx) error {
	user, _ := auth.GetUserData(c)

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
		[2]interface{}{"showAllTypesLink", true},
	)
}

// ShowProductsByType displays all products of a specific type
func ShowProductsByType(c *fiber.Ctx) error {
	productType := c.Params("type")
	user, _ := auth.GetUserData(c)

	// Validate product type
	validTypes := map[string]bool{
		"audio":   true,
		"midi":    true,
		"samples": true,
	}

	if !validTypes[productType] {
		return c.Status(404).SendString("Невідомий тип продукту")
	}

	// Get all products of the specified type
	products, err := productManagement.GetAllProductsByType(productType)
	if err != nil {
		fmt.Println("Error fetching products:", err)
		return c.Status(500).SendString("Помилка завантаження продуктів")
	}

	return render.RenderTemplate(c, "catalogue_type.html",
		[2]interface{}{"user", user},
		[2]interface{}{"products", products},
		[2]interface{}{"type", productType},
		[2]interface{}{"title", formatTypeTitle(productType)},
	)
}

// ShowVendorProducts displays all products from a specific vendor
func ShowVendorProducts(c *fiber.Ctx) error {
	vendorUsername := c.Params("username")
	user, _ := auth.GetUserData(c)

	// Get all products by vendor
	productsByType, vendorInfo, err := productManagement.GetAllProductsByVendor(vendorUsername)
	if err != nil {
		fmt.Println("Error fetching vendor products:", err)
		return c.Status(500).SendString("Помилка завантаження продуктів")
	}

	return render.RenderTemplate(c, "vendor_products.html",
		[2]interface{}{"user", user},
		[2]interface{}{"productsByType", productsByType},
		[2]interface{}{"vendor", vendorInfo},
	)
}

// formatTypeTitle returns a human-readable title for a product type
func formatTypeTitle(productType string) string {
	switch productType {
	case "audio":
		return "Аудіо"
	case "midi":
		return "MIDI"
	case "samples":
		return "Семпли"
	default:
		return productType
	}
}
