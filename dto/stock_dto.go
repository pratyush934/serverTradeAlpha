package dto

type StockDTO struct {
	Name   string  `json:"name"`
	Sector string  `json:"sector"`
	Price  float64 `json:"price"`
}
