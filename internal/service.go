// Оновлення internal/service.go для включення обробників плейлистів
package internal

import (
	"fmt"

	"github.com/DjentBoiiii/marketplace/internal/admin"
	"github.com/DjentBoiiii/marketplace/internal/auth"
	"github.com/DjentBoiiii/marketplace/internal/cart"
	"github.com/DjentBoiiii/marketplace/internal/catalogue"
	"github.com/DjentBoiiii/marketplace/internal/comments"
	"github.com/DjentBoiiii/marketplace/internal/filetransfer"
	"github.com/DjentBoiiii/marketplace/internal/likes"
	"github.com/DjentBoiiii/marketplace/internal/player"
	"github.com/DjentBoiiii/marketplace/internal/playlist"
	"github.com/DjentBoiiii/marketplace/internal/productManagement"
	"github.com/gofiber/fiber/v2"
)

func SetupHandlers(app *fiber.App) {
	app.Get("/", index)
	auth.SetupAuthHandlers(app)
	filetransfer.SetupUploadHandlers(app)
	productManagement.SetupProductHandlers(app)
	cart.SetupCartHandlers(app)
	playlist.SetupPlaylistHandlers(app)
	player.RegisterRoutes(app)
	comments.SetupCommentRoutes(app)
	likes.SetupLikeRoutes(app)
	filetransfer.SetupUploadHandlers(app)
	catalogue.SetupCatalogueRoutes(app)
	admin.SetupAdminHandlers(app)
}

func StartService() {
	app := fiber.New(fiber.Config{
		BodyLimit: 50 * 1024 * 1024,
	})
	app.Static("/static", "/marketplace/web/static")
	fmt.Println("Service starts on port 69420")
	SetupHandlers(app)
	app.Listen(":3000")
}
