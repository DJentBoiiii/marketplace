package models

import (
	"time"
)

type Comment struct {
	ID           int       `json:"id"`
	UserID       int       `json:"user_id"`
	Username     string    `json:"username"`
	ProductID    int       `json:"product_id"`
	Comment      string    `json:"comment"`
	LikesProduct bool      `json:"likes_product"`
	CreatedAt    time.Time `json:"created_at"`
	ProfilePhoto string    `json:"profile_photo"`
}
