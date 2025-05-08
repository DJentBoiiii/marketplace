package admin

import (
	"github.com/DjentBoiiii/marketplace/internal/auth"
	"github.com/DjentBoiiii/marketplace/internal/models"
	"github.com/DjentBoiiii/marketplace/internal/render"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gofiber/fiber/v2"
)

func ListUsers(c *fiber.Ctx) error {
	user, err := auth.GetUserData(c)
	if err != nil {
		return c.Status(401).SendString("Необхідно увійти в систему")
	}

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

	data := render.TemplateData{
		"users": users,
		"user":  user,
	}

	return render.RenderTemplate(c, "admin_users.html", data)
}
