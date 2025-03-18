package models

type CartItem struct {
	ProductID   int
	Name        string
	Price       int
	Type        string
	Description string
	Owner       string
	ImageURL    string
	Genre       string
}
