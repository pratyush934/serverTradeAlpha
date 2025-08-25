package models

import "time"

type PortFolioStock struct {
	Id           string    `gorm:"primaryKey; type:varchar(151)" json:"id"`
	StockId      string    `gorm:"not null" json:"stockId"`
	PortFolioId  string    `gorm:"not null" json:"portFolioId"`
	Quantity     int       `gorm:"default:0" json:"quantity"`
	AveragePrice float64   `gorm:"default:0" json:"averagePrice"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}
