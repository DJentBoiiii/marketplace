package auth

import (
	"database/sql"
	"fmt"

	"github.com/DjentBoiiii/marketplace/internal/models"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

func GetUserFromDB(username string) (*models.Account, error) {
	DB, err := sql.Open("mysql", DB_USER+":"+DB_PASSWORD+"@tcp(boku-no-sukele:3306)/"+DB_NAME)
	if err != nil {
		return nil, fmt.Errorf("помилка підключення до БД: %w", err)
	}
	defer DB.Close()

	var user models.Account
	// Додаємо поля bio та profile_pic
	err = DB.QueryRow(`
	    SELECT id, username, email, is_admin, 
	    COALESCE(bio, '') as bio, 
	    COALESCE(profile_photo, '') as profile_photo 
	    FROM Users WHERE username = ?`, username).
		Scan(&user.Id, &user.Username, &user.Email, &user.Is_admin, &user.Bio, &user.Profile_pic)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("користувача не знайдено")
		}
		return nil, fmt.Errorf("помилка отримання даних користувача: %w", err)
	}
	return &user, nil
}

func GetUserData(c *fiber.Ctx) (models.Account, error) {
	var user models.Account
	cookie := c.Cookies("jwt")
	if cookie == "" {
		user.Fill_default()
		return user, nil
	}
	token, err := jwt.Parse(cookie, func(token *jwt.Token) (interface{}, error) {
		return []byte(JWT_SECRET), nil
	})
	if err != nil || !token.Valid {
		return user, fmt.Errorf("invalid token")
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return user, fmt.Errorf("invalid token claims")
	}

	user.Id = int(claims["user_id"].(float64))
	user.Username = claims["username"].(string)
	user.Email = claims["email"].(string)
	user.Is_admin = claims["is_admin"].(bool)
	return user, nil
}
