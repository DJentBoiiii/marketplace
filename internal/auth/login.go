package auth

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/DjentBoiiii/marketplace/internal/render"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

func login(c *fiber.Ctx) error {
	fmt.Println("Login page visited")
	return render.RenderTemplate(c, "login.html")
}

func processLogin(c *fiber.Ctx) error {
	username := c.FormValue("username")
	password := c.FormValue("password")
	DB, _ = sql.Open("mysql", DB_USER+":"+DB_PASSWORD+"@tcp(boku-no-sukele:3306)/"+DB_NAME)
	if username == "" || password == "" {
		return c.Status(400).SendString("Заповніть всі поля")
	}

	var id int
	var dbPassword string
	var dbUsername string
	var dbEmail string
	var IsAdmin bool
	err := DB.QueryRow("SELECT * FROM Users WHERE name = ?", username).Scan(&id, &dbUsername, &dbEmail, &dbPassword, &IsAdmin)
	if err != nil {
		return c.Status(500).SendString("Помилка входу")
	}
	hashed_password := hash_pwd(password)

	if dbPassword != hashed_password {
		return c.Status(400).SendString("Неправильний пароль")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":  id,
		"username": username,
		"email":    dbEmail,
		"is_admin": false,
		"exp":      time.Now().Add(time.Hour * 72).Unix(),
	})

	tokenString, err := token.SignedString([]byte(JWT_SECRET))
	if err != nil {
		return c.Status(500).SendString("Помилка генерації токена")
	}

	c.Cookie(&fiber.Cookie{
		Name:  "jwt",
		Value: tokenString,
	})
	fmt.Println("User logged in")
	return c.Redirect("/")
}

func Logout(c *fiber.Ctx) error {
	c.ClearCookie("jwt")
	return c.Redirect("/")
}

func IsLoggedIn(c *fiber.Ctx) bool {
	cookie := c.Cookies("jwt")
	if cookie == "" {
		return false
	}

	token, err := jwt.Parse(cookie, func(token *jwt.Token) (interface{}, error) {
		return []byte(JWT_SECRET), nil
	})

	if err != nil || !token.Valid {
		return false
	}

	return true
}

func LoginRequired() func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		if !IsLoggedIn(c) {
			return c.Redirect("/login")
		}
		return c.Next()
	}
}
