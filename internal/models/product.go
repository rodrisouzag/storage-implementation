package models

type Product struct {
	ID       int     `json:"id"`
	Name     string  `json:"nombre"`
	Category string  `json:"tipo"`
	Count    int     `json:"cantidad"`
	Price    float64 `json:"precio"`
}
