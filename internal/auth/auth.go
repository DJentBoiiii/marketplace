package auth

import (
	"crypto/sha256"
	"fmt"

	"github.com/DjentBoiiii/marketplace/config"
	"github.com/DjentBoiiii/marketplace/internal/db"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gofiber/fiber/v2"
)

var (
	err error
	DB  = db.DB
)

func hash_pwd(password string) string {
	passwordWithKey := password + config.SHA_SECRET
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
	app.Get("/profile/edit/:username", LoginRequired(), EditProfile)
	app.Post("/profile/update/:username", LoginRequired(), UpdateProfile)
	app.Post("/profile/change-password/:username", LoginRequired(), ChangePassword)
}
