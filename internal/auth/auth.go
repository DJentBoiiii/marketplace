package auth

import (
	"crypto/sha256"
	"database/sql"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

var (
	_           = godotenv.Load("/marketplace/.env")
	DB_USER     = os.Getenv("MYSQL_USER")
	DB_PASSWORD = os.Getenv("MYSQL_PASSWORD")
	DB_NAME     = os.Getenv("MYSQL_DATABASE")
	JWT_SECRET  = os.Getenv("JWT_SECRET")
	SHA_SECRET  = os.Getenv("SHA_SECRET")
	DB          *sql.DB
)

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
