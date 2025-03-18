package models

import "time"

type Product struct {
	Name        string
	Price       int
	Type        string
	Description string
	Owner       string
	ImageURL    string
	CreatedAt   time.Time
}
