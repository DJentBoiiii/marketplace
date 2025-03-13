package handlers

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	"github.com/DjentBoiiii/marketplace/internal"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/joho/godotenv"
)

var _ = godotenv.Load("/marketplace/.env")
var DB_USER = os.Getenv("MYSQL_USER")
var DB_PASSWORD = os.Getenv("MYSQL_PASSWORD")
var DB_NAME = os.Getenv("MYSQL_DATABASE")
var JWT_SECRET = os.Getenv("JWT_SECRET")
var DB *sql.DB

func register(c *fiber.Ctx) error {
	fmt.Println("Register page visited")
	return internal.RenderTemplate(c, "register.html", nil)
}

func login(c *fiber.Ctx) error {
	fmt.Println("Login page visited")
	return internal.RenderTemplate(c, "login.html", nil)
}

func ProcessRegister(c *fiber.Ctx) error {
	var err error
	DB, err = sql.Open("mysql", DB_USER+":"+DB_PASSWORD+"@tcp(boku-no-sukele:3306)/"+DB_NAME)
	if err != nil {
		fmt.Println(err)
		return c.Status(500).SendString("Помилка підключення до бази даних:" + err.Error())
	}

	username := c.FormValue("username")
	password := c.FormValue("password")
	email := c.FormValue("email")

	if username == "" || password == "" || email == "" {
		return c.Status(400).SendString("Заповніть всі поля")
	}

	_, err = DB.Exec("INSERT INTO Users (name, email, password) VALUES (?, ?, ?)", username, email, password)
	fmt.Println(err)
	if err != nil {
		return c.Status(500).SendString("Помилка реєстрації")
	}
	fmt.Println("User registered")
	return c.Redirect("/")
}

func ProcessLogin(c *fiber.Ctx) error {
	username := c.FormValue("username")
	password := c.FormValue("password")
	DB, _ = sql.Open("mysql", DB_USER+":"+DB_PASSWORD+"@tcp(boku-no-sukele:3306)/"+DB_NAME)
	if username == "" || password == "" {
		return c.Status(400).SendString("Заповніть всі поля")
	}

	var id int
	var dbPassword string
	err := DB.QueryRow("SELECT id, password FROM Users WHERE name = ?", username).Scan(&id, &dbPassword)
	if err != nil {
		return c.Status(500).SendString("Помилка входу")
	}

	if dbPassword != password {
		return c.Status(400).SendString("Неправильний пароль")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": id,
		"exp":     time.Now().Add(time.Hour * 72).Unix(),
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

func SetupAuthHandlers(app *fiber.App) {
	app.Get("/register", register)
	app.Post("/submit", ProcessRegister)
	app.Get("/login", login)
	app.Post("/login", ProcessLogin)
	app.Get("/logout", Logout)

}
