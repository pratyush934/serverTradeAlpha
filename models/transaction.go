package models

import "time"

type TransactionModel struct {
	Id          string    `gorm:"primaryKey;type:varchar(151)" json:"id"`
	UserId      string    `gorm:"not null" json:"userId"`
	PortFolioId string    `gorm:"not null" json:"portFolioId"`
	StockId     string    `gorm:"not null" json:"stockId"`
	Quantity    int       `gorm:"default:0" json:"quantity"`
	Price       float64   `gorm:"default:0" json:"price"`
	Type        string    `json:"type"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}
