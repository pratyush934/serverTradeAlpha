package models

import "time"

type WatchListStockModel struct {
	Id          string `gorm:"primaryKey" json:"id"`
	WatchListId string `gorm:"not null" json:"watchListId"`
	StockId     string `gorm:"not null" json:"stockId"`
	CreatedAt   time.Time
}
