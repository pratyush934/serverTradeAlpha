package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/pratyush934/tradealpha/server/database"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

type WatchListStockModel struct {
	Id          string `gorm:"primaryKey" json:"id"`
	WatchListId string `gorm:"not null" json:"watchListId"`
	StockId     string `gorm:"not null" json:"stockId"`
	CreatedAt   time.Time
}

func (w *WatchListStockModel) BeforeCreate(tx *gorm.DB) error {
	w.Id = uuid.New().String()
	w.CreatedAt = time.Now()

	return nil
}

func (w *WatchListStockModel) CreateWatchListStockModel() (*WatchListStockModel, error) {
	if err := database.DB.Create(w).Error; err != nil {
		log.Error().Err(err).Msg("issue exist in watchlist_stock_model/CreateWatchListStockModel")
		return nil, err
	}
	return w, nil
}

func DeleteWatchListStockModel(id string) error {
	return database.DB.Where("id = ?", id).Delete(&WatchListStockModel{}).Error
}

func GetWatchListStockModelByStockId(stockId string) (*WatchListStockModel, error) {
	var watchListStockModel WatchListStockModel
	if err := database.DB.Where("stock_id = ?", stockId).First(&watchListStockModel).Error; err != nil {
		log.Error().Err(err).Msg("Issue exist in the watchlist_stock_model/GetWatchListStockModel")
		return nil, err
	}
	return &watchListStockModel, nil
}
func GetWatchListStockModelById(id string) (*WatchListStockModel, error) {
	var watchListStockModel WatchListStockModel
	if err := database.DB.Where("id = ?", id).First(&watchListStockModel).Error; err != nil {
		log.Error().Err(err).Msg("Issue exist in the watchlist_stock_model/GetWatchListStockModel")
		return nil, err
	}
	return &watchListStockModel, nil
}
