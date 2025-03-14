package auth

import (
	"crypto/sha256"
	"database/sql"
	"fmt"
	"os"
	"time"

	"github.com/DjentBoiiii/marketplace/internal/models"
	"github.com/DjentBoiiii/marketplace/internal/render"
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
var SHA_SECRET = os.Getenv("SHA_SECRET")
var DB *sql.DB

func register(c *fiber.Ctx) error {
	fmt.Println("Register page visited")
	return render.RenderTemplate(c, "register.html")
}

func login(c *fiber.Ctx) error {
	fmt.Println("Login page visited")
	return render.RenderTemplate(c, "login.html")
}

func processRegister(c *fiber.Ctx) error {
	var err error
	DB, err = sql.Open("mysql", DB_USER+":"+DB_PASSWORD+"@tcp(boku-no-sukele:3306)/"+DB_NAME)
	if err != nil {
		fmt.Println(err)
		return c.Status(500).SendString("Помилка підключення до бази даних:" + err.Error())
	}

	username := c.FormValue("username")
	password := c.FormValue("password")
	email := c.FormValue("email")

	hashed_pwd := hash_pwd(password)

	if username == "" || password == "" || email == "" {
		return c.Status(400).SendString("Заповніть всі поля")
	}

	_, err = DB.Exec("INSERT INTO Users (name, email, password) VALUES (?, ?, ?)", username, email, hashed_pwd)
	if err != nil {
		fmt.Println(err)
		return c.Status(500).SendString("Помилка реєстрації")
	}

	fmt.Println("User registered")
	return c.Redirect("/")
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

func GetUserFromDB(username string) (*models.Account, error) {
	DB, err := sql.Open("mysql", DB_USER+":"+DB_PASSWORD+"@tcp(boku-no-sukele:3306)/"+DB_NAME)
	if err != nil {
		return nil, fmt.Errorf("помилка підключення до БД: %w", err)
	}
	defer DB.Close()

	var user models.Account
	err = DB.QueryRow("SELECT id, name, email, is_admin FROM Users WHERE name = ?", username).
		Scan(&user.Id, &user.Username, &user.Email, &user.Is_admin)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("користувача не знайдено")
		}
		return nil, fmt.Errorf("помилка отримання даних користувача: %w", err)
	}
	return &user, nil
}

func Profile(c *fiber.Ctx) error {
	username := c.Params("username")
	user, err := GetUserFromDB(username)
	if err != nil {
		return c.Status(404).SendString("Користувач не знайдений")
	}

	return render.RenderTemplate(c, "profile.html",
		[2]interface{}{"user", user},
	)
}

func hash_pwd(password string) string {
	passwordWithKey := password + SHA_SECRET
	hash := sha256.New()
	hash.Write([]byte(passwordWithKey))
	hashedPassword := hash.Sum(nil)

	hashedPasswordHex := fmt.Sprintf("%x", hashedPassword)

	return hashedPasswordHex
}

func SetupAuthHandlers(app *fiber.App) {
	app.Get("/register", register)
	app.Post("/submit", processRegister)
	app.Get("/login", login)
	app.Post("/login", processLogin)
	app.Get("/logout", LoginRequired(), Logout)
	app.Get("/profile/:username", Profile)

}
