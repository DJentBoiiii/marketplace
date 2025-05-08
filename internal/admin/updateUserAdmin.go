package admin

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/gofiber/fiber/v2"
)

// UpdateUserAdmin toggles administrator privileges for a user
// Used to promote regular users to admins or demote admins to regular users
func UpdateUserAdmin(c *fiber.Ctx) error {
	userId := c.Params("id")
	adminStatus := c.FormValue("is_admin") == "true"

	_, err := DB.Exec("UPDATE Users SET is_admin = ? WHERE id = ?", adminStatus, userId)
	if err != nil {
		return c.Status(500).SendString("Помилка оновлення статусу адміністратора")
	}

	return c.Redirect("/admin/users")
}
