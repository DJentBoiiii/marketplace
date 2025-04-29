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
	DB, _ = sql.Open("mysql", DB_USER+":"+DB_PASSWORD+"@tcp("+DB_HOST+":3306)/"+DB_NAME)
	if username == "" || password == "" {
		return c.Status(400).SendString("Заповніть всі поля")
	}

	var id int
	var dbPassword string
	var dbUsername string
	var dbEmail string
	var IsAdmin bool
	var profilePhoto string
	var bio string
	var createdAt []uint8
	err := DB.QueryRow(`
	SELECT id, username, email, password, is_admin, COALESCE(profile_photo, ''), COALESCE(bio, ''), created_at 
	FROM Users WHERE username = ?`, username).Scan(
		&id, &dbUsername, &dbEmail, &dbPassword, &IsAdmin, &profilePhoto, &bio, &createdAt)
	if err == nil {
		createdAtTime, err := time.Parse("2006-01-02 15:04:05", string(createdAt))
		if err != nil {
			return c.Status(500).SendString("Помилка обробки дати")
		}
		fmt.Println("User created at:", createdAtTime)
	}
	if err != nil {
		fmt.Println(err)
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
		"is_admin": IsAdmin,
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
