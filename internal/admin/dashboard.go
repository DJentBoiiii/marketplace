package admin

import (
	"github.com/DjentBoiiii/marketplace/internal/auth"
	"github.com/DjentBoiiii/marketplace/internal/models"
	"github.com/DjentBoiiii/marketplace/internal/render"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gofiber/fiber/v2"
)

func Dashboard(c *fiber.Ctx) error {

	user, err := auth.GetUserData(c)
	if err != nil {
		return c.Status(401).SendString("Необхідно увійти в систему")
	}
	var stats models.AdminStats

	err = DB.QueryRow("SELECT COUNT(*) FROM Users").Scan(&stats.UserCount)
	if err != nil {
		return c.Status(500).SendString("Помилка отримання статистики користувачів")
	}

	err = DB.QueryRow("SELECT COUNT(*) FROM Products").Scan(&stats.ProductCount)
	if err != nil {
		return c.Status(500).SendString("Помилка отримання статистики продуктів")
	}

	err = DB.QueryRow("SELECT COUNT(*) FROM Purchases").Scan(&stats.PurchaseCount)
	if err != nil {
		return c.Status(500).SendString("Помилка отримання статистики покупок")
	}

	rows, err := DB.Query(`
		SELECT id, username, email, created_at 
		FROM Users 
		ORDER BY created_at DESC LIMIT 5`)
	if err != nil {
		return c.Status(500).SendString("Помилка отримання останніх користувачів")
	}
	defer rows.Close()

	var latestUsers []models.UserInfo
	for rows.Next() {
		var user models.UserInfo
		var createdAt []uint8
		err := rows.Scan(&user.Id, &user.Username, &user.Email, &createdAt)
		if err != nil {
			continue
		}
		user.CreatedAt = string(createdAt)
		latestUsers = append(latestUsers, user)
	}

	return render.RenderTemplate(c, "admin_dashboard.html",
		[2]interface{}{"user", user},
		[2]interface{}{"stats", stats},
		[2]interface{}{"latestUsers", latestUsers},
		[2]interface{}{"active_page", "dashboard"},
	)
}
