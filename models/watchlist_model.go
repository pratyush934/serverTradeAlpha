package models

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/pratyush934/tradealpha/server/database"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

type WatchListModel struct {
	Id             string                `gorm:"primaryKey;type:varchar(151)" json:"id"`
	UserId         string                `gorm:"not null" json:"userId"`
	Name           string                `gorm:"not null" json:"name"`
	Description    string                `gorm:"not null;size:1000" json:"description"`
	WatchListStock []WatchListStockModel `gorm:"foreignKey:WatchListId" json:"watchListStock"`
	CreatedAt      time.Time             `json:"createdAt"`
	UpdatedAt      time.Time             `json:"updatedAt"`
}

func (w *WatchListModel) BeforeCreate(tx *gorm.DB) error {
	w.Id = uuid.New().String()
	w.CreatedAt = time.Now()
	w.UpdatedAt = time.Now()

	return nil
}

func (w *WatchListModel) Create() (*WatchListModel, error) {
	if err := database.DB.Create(w).Error; err != nil {
		log.Error().Err(err).Msg("issue persist in the watchlist_model/Create")
		return nil, err
	}
	return w, nil
}

func GetWatchListById(id string) (*WatchListModel, error) {
	var watchList WatchListModel
	if err := database.DB.Preload("WatchListStock").Where("id = ?", id).First(&watchList).Error; err != nil {
		log.Error().Err(err).Msg("issue persist in the watchlist_model/GetWatchListById")
		return nil, err
	}
	return &watchList, nil
}

func GetWatchListsByUserId(userId string) (*[]WatchListModel, error) {
	var watchList []WatchListModel
	if err := database.DB.Where("user_id = ?", userId).Find(&watchList).Error; err != nil {
		log.Error().Err(err).Msg("issue persist in the watchlist_model/GetWatchListsByUserId")
		return nil, err
	}
	return &watchList, nil
}

func AddStockToWatchlist(watchlistId, symbol string) error {
	// Verify watchlist exists
	_, err := GetWatchListById(watchlistId)
	if err != nil {
		log.Error().Err(err).Str("watchlist_id", watchlistId).Msg("Failed to find watchlist")
		return err
	}

	// Verify stock exists
	_, err = GetStockBySymbol(symbol)
	if err != nil {
		log.Error().Err(err).Str("symbol", symbol).Msg("Failed to find stock")
		return err
	}

	// Check if stock is already in watchlist
	var existing WatchListStockModel
	if err := database.DB.Where("watchlist_id = ? AND symbol = ?", watchlistId, symbol).First(&existing).Error; err == nil {
		log.Warn().Str("watchlist_id", watchlistId).Str("symbol", symbol).Msg("Stock already in watchlist")
		return fmt.Errorf("stock already in watchlist")
	}

	// Add stock to watchlist
	watchlistStock := WatchListStockModel{
		WatchListId: watchlistId,
		Symbol:      symbol,
		StockId:     existing.StockId,
	}
	if _, err := watchlistStock.CreateWatchListStockModel(); err != nil {
		log.Error().Err(err).Str("watchlist_id", watchlistId).Str("symbol", symbol).Msg("Failed to add stock to watchlist")
		return err
	}

	log.Info().Str("watchlist_id", watchlistId).Str("symbol", symbol).Msg("Stock added to watchlist")
	return nil
}

func RemoveStockFromWatchlist(watchlistId, symbol string) error {
	// Verify watchlist exists
	if _, err := GetWatchListById(watchlistId); err != nil {
		log.Error().Err(err).Str("watchlist_id", watchlistId).Msg("Failed to find watchlist")
		return err
	}

	// Delete stock from watchlist
	if err := database.DB.Where("watchlist_id = ? AND symbol = ?", watchlistId, symbol).Delete(&WatchListStockModel{}).Error; err != nil {
		log.Error().Err(err).Str("watchlist_id", watchlistId).Str("symbol", symbol).Msg("Failed to remove stock from watchlist")
		return err
	}

	log.Info().Str("watchlist_id", watchlistId).Str("symbol", symbol).Msg("Stock removed from watchlist")

	return nil
}

func DeleteWatchList(id string) error {
	return database.DB.Where("id = ?", id).Delete(&WatchListModel{}).Error
}

func GetAllWatchListStocksByUserId(userId string) ([]WatchListStockModel, error) {

	var watchLists []WatchListModel
	if err := database.DB.Where("user_id = ?", userId).Find(&watchLists).Error; err != nil {
		return nil, err
	}

	var allStocks []WatchListStockModel

	for _, wl := range watchLists {
		var stocks []WatchListStockModel
		if err := database.DB.Where(&WatchListStockModel{WatchListId: wl.Id}).Find(&stocks).Error; err != nil {
			log.Error().Err(err).Msg("issue persist in the GetAllWatchListStocksByUserId")
			return nil, err
		}
		allStocks = append(allStocks, stocks...)
	}
	return allStocks, nil
}

func UpdateWatchlist(w WatchListModel) error {
	return database.DB.Updates(&w).Error
}
