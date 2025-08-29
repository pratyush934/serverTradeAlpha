package dto

type TransactionDTO struct {
	Quantity int     `json:"quantity"`
	Price    float64 `json:"price"`
	Status   string  `json:"status"`
	Type     string  `json:"type"`
}
