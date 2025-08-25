package models

import "time"

type Stock struct {
	Id             string             `gorm:"primaryKey;type:varchar(151)" json:"id"`
	Name           string             `gorm:"not null" json:"name"`
	Sector         string             `json:"sector"`
	Price          float64            `gorm:"not null" json:"price"`
	PortFolioStock []PortFolioStock   `gorm:"foreignKey:stockId" json:"portFolioStock"`
	OrderBook      []OrderBookModel   `gorm:"foreignKey:stockId" json:"orderBook"`
	Transaction    []TransactionModel `gorm:"foreignKey:stockId" json:"transaction"`
	CreatedAt      time.Time          `json:"createdAt"`
	UpdatedAt      time.Time          `json:"updatedAt"`
}
