package catalogue

import (
	"github.com/gofiber/fiber/v2"
)

func SetupCatalogueRoutes(app *fiber.App) {
	app.Get("/catalogue", ShowCatalogue)
	app.Get("/catalogue/:type", ShowProductsByType)
	app.Get("/test/vendor/:username", ShowVendorProducts)
}
