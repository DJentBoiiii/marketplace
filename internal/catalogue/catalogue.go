package catalogue

import (
	"github.com/DjentBoiiii/marketplace/internal/db"
	"github.com/gofiber/fiber/v2"
)

var (
	DB = db.DB
)

func SetupCatalogueRoutes(app *fiber.App) {
	app.Get("/catalogue", ShowCatalogue)
	app.Get("/catalogue/:type", ShowProductsByType)
	app.Get("/test/vendor/:username", ShowVendorProducts)
	app.Get("/artists", ShowArtists) // Add route for artists catalogue
}
