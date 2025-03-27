package models

type PlaylistItem struct {
	ItemId      int
	ProductID   int
	Name        string
	Price       int
	Type        string
	Description string
	Owner       string
	ImageURL    string
	Genre       string
}
