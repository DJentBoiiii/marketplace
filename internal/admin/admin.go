package admin

import (
	"github.com/DjentBoiiii/marketplace/internal/db"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gofiber/fiber/v2"
)

var DB = db.DB

func SetupAdminHandlers(app *fiber.App) {
	adminGroup := app.Group("/admin", AdminRequired())
	adminGroup.Get("/", Dashboard)
	adminGroup.Get("/users", ListUsers)
	adminGroup.Post("/users/:id/toggle-admin", UpdateUserAdmin)
	adminGroup.Post("/users/:id/delete", DeleteUser)

	adminGroup.Get("/products", ListProducts)
	adminGroup.Post("/products/:id/delete", DeleteProduct)
}
