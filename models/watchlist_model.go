package models

import (
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

/*
CreateWatchlist: Creates a new watchlist for a user.
GetWatchlistById: Retrieves a watchlist by ID, preloading stocks.
GetWatchlistsByUserId: Fetches all watchlists for a user.
DeleteWatchlist: Deletes a watchlist by ID.
AddStockToWatchlist: Adds a stock to a watchlist.
RemoveStockFromWatchlist: Removes a stock from a watchlist.
GetWatchlistStocks: Retrieves all stocks in a watchlist.
*/

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
	if err := database.DB.Where("id = ?", id).First(&watchList).Error; err != nil {
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
