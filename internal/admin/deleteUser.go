package admin

import (
	"database/sql"
	"fmt"

	"github.com/DjentBoiiii/marketplace/internal/auth"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gofiber/fiber/v2"
)

// DeleteUser removes a user from the system
// Prevents self-deletion and removes all user data from the database
func DeleteUser(c *fiber.Ctx) error {
	userId := c.Params("id")

	DB, err := sql.Open("mysql", auth.DB_USER+":"+auth.DB_PASSWORD+"@tcp("+auth.DB_HOST+":3306)/"+auth.DB_NAME)
	if err != nil {
		return c.Status(500).SendString("Помилка підключення до БД")
	}
	defer DB.Close()

	currentUser, _ := auth.GetUserData(c)
	if fmt.Sprintf("%d", currentUser.Id) == userId {
		return c.Status(400).SendString("Ви не можете видалити свій власний обліковий запис")
	}

	_, err = DB.Exec("DELETE FROM Users WHERE id = ?", userId)
	if err != nil {
		return c.Status(500).SendString("Помилка видалення користувача")
	}

	return c.Redirect("/admin/users")
}
