package models

import (
	"time"
)

type ProductLike struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	Username  string    `json:"username,omitempty"`
	ProductID int       `json:"product_id"`
	CreatedAt time.Time `json:"created_at"`
}

type ProductLikeResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
	IsLiked bool   `json:"is_liked"`
	Count   int    `json:"count,omitempty"`
}
