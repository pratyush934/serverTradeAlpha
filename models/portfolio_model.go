package models

import "time"

type PortFolio struct {
	Id             string             `gorm:"primaryKey;type:varchar(151)" json:"id"`
	UserId         string             `gorm:"not null" json:"userId"`
	Name           string             `gorm:"not null" json:"name"`
	Title          string             `gorm:"not null" json:"title"`
	TotalValue     float64            `gorm:"default:0" json:"totalValue"`
	Description    string             `gorm:"not null" json:"description"`
	Transaction    []TransactionModel `gorm:"foreignKey:portFolioId" json:"transaction"`
	PortFolioStock []PortFolioStock   `gorm:"foreignKey:portFolioId" json:"portFolioStock"`
	CreatedAt      time.Time          `json:"createdAt"`
	UpdatedAt      time.Time          `json:"updatedAt"`
}
