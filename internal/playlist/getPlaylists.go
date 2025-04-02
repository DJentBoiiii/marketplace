// internal/playlist/getPlaylists.go
package playlist

import (
	"github.com/DjentBoiiii/marketplace/internal/models"
	_ "github.com/go-sql-driver/mysql"
)

func GetUserPlaylists(userId int) ([]models.Playlist, error) {
	rows, err := DB.Query(`
		SELECT id, name, created_at
		FROM Playlists
		WHERE user_id = ?
		ORDER BY created_at DESC`, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var playlists []models.Playlist
	for rows.Next() {
		var p models.Playlist
		var createdAtBytes []uint8
		if err := rows.Scan(&p.Id, &p.Name, &createdAtBytes); err != nil {
			return nil, err
		}

		// Конвертуємо рядок в time.Time
		p.CreatedAt, err = parseCreatedTime(string(createdAtBytes))
		if err != nil {
			return nil, err
		}

		// Отримуємо кількість треків
		err = DB.QueryRow(`
			SELECT COUNT(*) 
			FROM PlaylistItems 
			WHERE playlist_id = ?`, p.Id).Scan(&p.ItemCount)
		if err != nil {
			return nil, err
		}

		playlists = append(playlists, p)
	}
	return playlists, nil
}

func GetPlaylistItems(playlistId string) ([]models.PlaylistItem, error) {
	rows, err := DB.Query(`
		SELECT pi.id AS item_id, p.id, p.name, p.price, p.type, 
			   p.description, u.username AS owner, p.image_url, p.genre
		FROM PlaylistItems pi
		JOIN Products p ON pi.product_id = p.id
		JOIN Users u ON p.vendor = u.username
		WHERE pi.playlist_id = ?`, playlistId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []models.PlaylistItem
	for rows.Next() {
		var item models.PlaylistItem
		var imagePath string
		if err := rows.Scan(
			&item.ItemId, &item.ProductID, &item.Name, &item.Price, &item.Type,
			&item.Description, &item.Owner, &imagePath, &item.Genre); err != nil {
			return nil, err
		}
		item.ImageURL = "/" + imagePath
		items = append(items, item)
	}
	return items, nil
}

func parseCreatedTime(timeStr string) (string, error) {
	return timeStr, nil
}
