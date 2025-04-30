package admin

import (
	"database/sql"

	"github.com/DjentBoiiii/marketplace/internal/auth"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gofiber/fiber/v2"
)

// UpdateUserAdmin toggles administrator privileges for a user
// Used to promote regular users to admins or demote admins to regular users
func UpdateUserAdmin(c *fiber.Ctx) error {
	userId := c.Params("id")
	adminStatus := c.FormValue("is_admin") == "true"

	DB, err := sql.Open("mysql", auth.DB_USER+":"+auth.DB_PASSWORD+"@tcp("+auth.DB_HOST+":3306)/"+auth.DB_NAME)
	if err != nil {
		return c.Status(500).SendString("Помилка підключення до БД")
	}
	defer DB.Close()

	_, err = DB.Exec("UPDATE Users SET is_admin = ? WHERE id = ?", adminStatus, userId)
	if err != nil {
		return c.Status(500).SendString("Помилка оновлення статусу адміністратора")
	}

	return c.Redirect("/admin/users")
}
