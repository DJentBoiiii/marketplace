package catalogue

import (
	"log"

	"github.com/DjentBoiiii/marketplace/internal/auth"
	"github.com/DjentBoiiii/marketplace/internal/render"
	"github.com/gofiber/fiber/v2"
)

// Artist represents a user who is an artist
type Artist struct {
	Username     string
	ProfilePic   string
	Bio          string
	ProductCount int
}

// ShowArtists handles the display of all artists
func ShowArtists(c *fiber.Ctx) error {
	// Get current user data if logged in
	user, _ := auth.GetUserData(c)

	// Query to get all artists (users with is_artist flag = true)
	rows, err := DB.Query(`
		SELECT u.username, COALESCE(u.profile_photo, ''), COALESCE(u.bio, ''), COUNT(p.id) as product_count
		FROM Users u
		LEFT JOIN Products p ON u.username = p.vendor
		WHERE u.is_artist = TRUE
		GROUP BY u.username, u.profile_photo, u.bio
		ORDER BY product_count DESC
	`)
	if err != nil {
		log.Printf("Error querying artists: %v", err)
		return c.Status(500).SendString("Database error")
	}
	defer rows.Close()

	// Process the results
	artists := []Artist{}
	for rows.Next() {
		var artist Artist
		if err := rows.Scan(&artist.Username, &artist.ProfilePic, &artist.Bio, &artist.ProductCount); err != nil {
			log.Printf("Error scanning artist row: %v", err)
			continue
		}

		// If profile pic is empty, use a default image
		if artist.ProfilePic == "" {
			artist.ProfilePic = "static/images/default_profile.jpg"
		}

		artists = append(artists, artist)
	}

	// Check for errors from iterating over rows
	if err = rows.Err(); err != nil {
		log.Printf("Error iterating over rows: %v", err)
		return c.Status(500).SendString("Database error")
	}

	data := render.TemplateData{
		"Title":   "Artists - DSA Marketplace",
		"User":    user,
		"Artists": artists,
	}

	// Render the template with artists data using render.RenderTemplate
	return render.RenderTemplate(c, "catalogue_artists.html", data)
}
