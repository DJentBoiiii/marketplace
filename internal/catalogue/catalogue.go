package catalogue

import (
	"github.com/gofiber/fiber/v2"
)

func SetupCatalogueRoutes(app *fiber.App) {
	app.Get("/catalogue", ShowCatalogue)
	app.Get("/catalogue/:type", ShowProductsByType)
	app.Get("/test/vendor/:username", ShowVendorProducts)
	app.Get("/artists", ShowArtists) // Add route for artists catalogue

	// Add a redirect from old vendor route to new one
	app.Get("/catalogue/vendor/:username", func(c *fiber.Ctx) error {
		username := c.Params("username")
		return c.Redirect("/test/vendor/"+username, 301) // 301 is permanent redirect
	})
}
