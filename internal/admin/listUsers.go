package admin

import (
	"database/sql"

	"github.com/DjentBoiiii/marketplace/internal/auth"
	"github.com/DjentBoiiii/marketplace/internal/models"
	"github.com/DjentBoiiii/marketplace/internal/render"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gofiber/fiber/v2"
)

// ListUsers retrieves all registered users and renders them in admin panel
// Used in the admin users management page to display and manage all system users
func ListUsers(c *fiber.Ctx) error {
	user, err := auth.GetUserData(c)
	if err != nil {
		return c.Status(401).SendString("Необхідно увійти в систему")
	}

	DB, err := sql.Open("mysql", auth.DB_USER+":"+auth.DB_PASSWORD+"@tcp("+auth.DB_HOST+":3306)/"+auth.DB_NAME)
	if err != nil {
		return c.Status(500).SendString("Помилка підключення до БД")
	}
	defer DB.Close()

	rows, err := DB.Query(`
		SELECT id, username, email, is_admin, created_at 
		FROM Users 
		ORDER BY created_at DESC`)
	if err != nil {
		return c.Status(500).SendString("Помилка отримання списку користувачів")
	}
	defer rows.Close()

	var users []models.UserInfo
	for rows.Next() {
		var userInfo models.UserInfo
		var createdAt []uint8
		err := rows.Scan(&userInfo.Id, &userInfo.Username, &userInfo.Email, &userInfo.IsAdmin, &createdAt)
		if err != nil {
			continue
		}
		userInfo.CreatedAt = string(createdAt)
		users = append(users, userInfo)
	}

	return render.RenderTemplate(c, "admin_users.html",
		[2]interface{}{"user", user},
		[2]interface{}{"users", users},
	)
}
