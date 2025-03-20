package auth

import (
	"database/sql"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/DjentBoiiii/marketplace/internal/models"
	"github.com/DjentBoiiii/marketplace/internal/render"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gofiber/fiber/v2"
)

func EditProfile(c *fiber.Ctx) error {
	// Get current user data
	user, err := GetUserData(c)
	if err != nil {
		return c.Status(401).SendString("Необхідно увійти в систему")
	}

	// Fetch the full user data from DB including bio and profile photo
	DB, err := sql.Open("mysql", DB_USER+":"+DB_PASSWORD+"@tcp(boku-no-sukele:3306)/"+DB_NAME)
	if err != nil {
		return c.Status(500).SendString("Помилка підключення до БД")
	}
	defer DB.Close()

	var bio, profilePhoto string
	err = DB.QueryRow("SELECT COALESCE(bio, ''), COALESCE(profile_photo, '') FROM Users WHERE id = ?", user.Id).
		Scan(&bio, &profilePhoto)
	if err != nil {
		return c.Status(500).SendString("Помилка отримання даних користувача")
	}

	// Add the additional fields to user struct
	user.Bio = bio
	user.Profile_pic = profilePhoto

	return render.RenderTemplate(c, "edit_profile.html",
		[2]interface{}{"user", user},
	)
}

func UpdateProfile(c *fiber.Ctx) error {
	// Get current user data
	user, err := GetUserData(c)
	if err != nil {
		return c.Status(401).SendString("Необхідно увійти в систему")
	}

	// Get form values
	email := c.FormValue("email")
	bio := c.FormValue("bio")

	// Connect to database
	DB, err := sql.Open("mysql", DB_USER+":"+DB_PASSWORD+"@tcp(boku-no-sukele:3306)/"+DB_NAME)
	if err != nil {
		return c.Status(500).SendString("Помилка підключення до БД")
	}
	defer DB.Close()

	var profilePhotoPath string = ""

	// Handle profile photo upload
	file, err := c.FormFile("profile_photo")
	if err == nil && file != nil {
		// Create directory if it doesn't exist
		uploadsDir := filepath.Join("..", "web", "static", "uploads", "profiles")
		if _, err := os.Stat(uploadsDir); os.IsNotExist(err) {
			os.MkdirAll(uploadsDir, 0755)
		}

		// Generate filename based on username to avoid collisions
		// Get file extension
		fileExt := filepath.Ext(file.Filename)
		// Create safe filename
		safeFilename := fmt.Sprintf("%s%s", user.Username, fileExt)
		// Full path
		filePath := filepath.Join(uploadsDir, safeFilename)

		// Save uploaded file
		src, err := file.Open()
		if err != nil {
			return c.Status(500).SendString("Помилка відкриття завантаженого файлу")
		}
		defer src.Close()

		// Create destination file
		dst, err := os.Create(filePath)
		if err != nil {
			return c.Status(500).SendString("Помилка створення файлу")
		}
		defer dst.Close()

		// Copy file contents
		if _, err = io.Copy(dst, src); err != nil {
			return c.Status(500).SendString("Помилка збереження файлу")
		}

		// Set profile photo path to be saved in DB
		profilePhotoPath = fmt.Sprintf("/static/uploads/profiles/%s", safeFilename)
	}

	// Update user profile in database
	var query string
	var args []interface{}

	if profilePhotoPath != "" {
		// If profile photo was uploaded, update everything
		query = "UPDATE Users SET email = ?, bio = ?, profile_photo = ? WHERE id = ?"
		args = []interface{}{email, bio, profilePhotoPath, user.Id}
	} else {
		// If no profile photo, just update email and bio
		query = "UPDATE Users SET email = ?, bio = ? WHERE id = ?"
		args = []interface{}{email, bio, user.Id}
	}

	_, err = DB.Exec(query, args...)
	if err != nil {
		return c.Status(500).SendString("Помилка оновлення профілю: " + err.Error())
	}

	// Redirect to profile page
	return c.Redirect(fmt.Sprintf("/profile/%s", user.Username))
}

func ChangePassword(c *fiber.Ctx) error {
	// Get current user data
	user, err := GetUserData(c)
	if err != nil {
		return c.Status(401).SendString("Необхідно увійти в систему")
	}

	// Get form values
	oldPassword := c.FormValue("old_password")
	newPassword := c.FormValue("new_password")
	confirmPassword := c.FormValue("confirm_password")

	// Validate input
	if oldPassword == "" || newPassword == "" || confirmPassword == "" {
		return c.Status(400).SendString("Заповніть всі поля")
	}

	if newPassword != confirmPassword {
		return c.Status(400).SendString("Паролі не співпадають")
	}

	// Connect to database
	DB, err := sql.Open("mysql", DB_USER+":"+DB_PASSWORD+"@tcp(boku-no-sukele:3306)/"+DB_NAME)
	if err != nil {
		return c.Status(500).SendString("Помилка підключення до БД")
	}
	defer DB.Close()

	// Verify current password
	var storedPassword string
	err = DB.QueryRow("SELECT password FROM Users WHERE id = ?", user.Id).Scan(&storedPassword)
	if err != nil {
		return c.Status(500).SendString("Помилка перевірки паролю")
	}

	hashedOldPassword := hash_pwd(oldPassword)
	if hashedOldPassword != storedPassword {
		return c.Status(400).SendString("Поточний пароль невірний")
	}

	// Hash and update new password
	hashedNewPassword := hash_pwd(newPassword)
	_, err = DB.Exec("UPDATE Users SET password = ? WHERE id = ?", hashedNewPassword, user.Id)
	if err != nil {
		return c.Status(500).SendString("Помилка оновлення паролю")
	}

	return c.Redirect(fmt.Sprintf("/profile/%s", user.Username))
}

func GetUserWithProfile(username string) (*models.Account, error) {
	DB, err := sql.Open("mysql", DB_USER+":"+DB_PASSWORD+"@tcp(boku-no-sukele:3306)/"+DB_NAME)
	if err != nil {
		return nil, fmt.Errorf("помилка підключення до БД: %w", err)
	}
	defer DB.Close()

	var user models.Account
	err = DB.QueryRow(`
		SELECT id, username, email, is_admin, COALESCE(bio, ''), COALESCE(profile_photo, '') 
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
