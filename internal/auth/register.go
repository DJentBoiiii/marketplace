package auth

import (
	"database/sql"
	"fmt"

	"github.com/DjentBoiiii/marketplace/internal/render"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gofiber/fiber/v2"
)

func register(c *fiber.Ctx) error {
	fmt.Println("Register page visited")
	return render.RenderTemplate(c, "register.html")
}

func processRegister(c *fiber.Ctx) error {
	var err error
	DB, err = sql.Open("mysql", DB_USER+":"+DB_PASSWORD+"@tcp("+DB_HOST+":3306)/"+DB_NAME)
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

	_, err = DB.Exec("INSERT INTO Users (username, email, password) VALUES (?, ?, ?)", username, email, hashed_pwd)
	if err != nil {
		fmt.Println(err)
		return c.Status(500).SendString("Помилка реєстрації")
	}

	fmt.Println("User registered")
	return c.Redirect("/")
}
