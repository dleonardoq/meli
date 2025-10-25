package models

type Product struct {
	ID             string         `json:"id"`
	Name           string         `json:"name"`
	ImageURL       string         `json:"image_url"`
	Description    string         `json:"description"`
	Price          float64        `json:"price"`
	Currency       string         `json:"currency"`
	Rating         float64        `json:"rating"`
	ReviewCount    int            `json:"review_count"`
	Specifications map[string]any `json:"specifications"`
	Category       string         `json:"category"`
	Brand          string         `json:"brand"`
	InStock        bool           `json:"in_stock"`
}
