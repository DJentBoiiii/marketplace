package models

import "time"

type Product struct {
	ID          int
	ItemID      int
	Name        string
	Price       int
	Type        string
	Description string
	Owner       string
	Genre       string
	ImageURL    string
	CreatedAt   time.Time
}
