package models

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/pratyush934/tradealpha/server/alphavantage"
	"github.com/pratyush934/tradealpha/server/database"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

type Stock struct {
	Id             string                `gorm:"primaryKey;type:varchar(151)" json:"id"`
	Name           string                `gorm:"not null" json:"name"`
	Sector         string                `json:"sector"`
	Price          float64               `gorm:"not null" json:"price"`
	Symbol         string                `json:"symbol"`
	WatchListStock []WatchListStockModel `gorm:"foreignKey:stockId" json:"watchListStock"`
	PortFolioStock []PortFolioStock      `gorm:"foreignKey:stockId" json:"portFolioStock"`
	Transaction    []TransactionModel    `gorm:"foreignKey:stockId" json:"transaction"`
	CreatedAt      time.Time             `json:"createdAt"`
	UpdatedAt      time.Time             `json:"updatedAt"`
}

func (s *Stock) BeforeCreate(tx *gorm.DB) error {
	s.Id = uuid.New().String()
	s.CreatedAt = time.Now()
	s.UpdatedAt = time.Now()

	return nil
}

func (s *Stock) BeforeUpdate(tx *gorm.DB) error {
	s.UpdatedAt = time.Now()
	return nil
}

func (s *Stock) CreateStock() (*Stock, error) {
	if err := database.DB.Create(s).Error; err != nil {
		log.Error().Err(err).Msg("issue persist in stock_model/CreateStock")
		return nil, err
	}
	return s, nil
}

func GetStockById(id string) (*Stock, error) {
	var stock Stock
	if err := database.DB.
		Where("stock_id = ?", id).
		Preload("WatchListStock").
		Preload("PortFolioStock").
		Preload("Transaction").
		First(&stock).Error; err != nil {
		log.Error().Err(err).Msg("issue persist in stock_model/GetStockById")
		return nil, err
	}
	return &stock, nil
}

func GetStockBySector(sector string, limit, offSet int) (*[]Stock, error) {
	var stock []Stock
	if err := database.DB.Where("sector = ?", sector).
		Preload("WatchListStock").
		Preload("PortFolioStock").
		Preload("Transaction").
		Limit(limit).Offset(offSet).
		Find(&stock).Error; err != nil {
		log.Error().Err(err).Msg("issue persist in stock_model/GetStockBySector")
		return nil, err
	}
	return &stock, nil
}

func GetAllStocks(limit, offset int) ([]Stock, error) {
	var stock []Stock
	if err := database.DB.
		Preload("WatchListStock").
		Preload("PortFolioStock").
		Preload("Transaction").
		Limit(limit).Offset(offset).
		Find(&stock).Error; err != nil {
		log.Error().Err(err).Msg("issue persist in stock_model/GetAllStocks")
		return nil, err
	}
	return stock, nil
}

func GetStockBySymbol(symbol string) (*Stock, error) {
	var stock Stock

	if err := database.DB.
		Where("symbol = ?", symbol).
		Preload("WatchListStock").
		Preload("PortFolioStock").
		Preload("Transaction").
		First(&stock).Error; err != nil {
		log.Error().Err(err).Msg("issue persist in the GetStockBySymbol")
	}

	return &stock, nil
}

func UpdateStock(stock Stock) error {
	return database.DB.Updates(&stock).Error
}

func FetchAndCacheStock(symbol string, logger *zerolog.Logger) (*Stock, error) {
	// Fetch stock overview (for Name and Sector)

	url := fmt.Sprintf("%s?function=OVERVIEW&symbol=%s&apikey=%s", alphavantage.AAlphaVantageBaseURL, symbol, os.Getenv("ALPHA_VANTAGE_KEY"))
	resp, err := http.Get(url)
	if err != nil {
		log.Error().Err(err).Str("symbol", symbol).Msg("Failed to fetch overview from Alpha Vantage")
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Error().Int("status", resp.StatusCode).Str("symbol", symbol).Msg("Alpha Vantage API returned non-200 status")
		return nil, fmt.Errorf("alpha Vantage API error")
	}

	var overview struct {
		Symbol string `json:"Symbol"`
		Name   string `json:"Name"`
		Sector string `json:"Sector"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&overview); err != nil {
		log.Error().Err(err).Str("symbol", symbol).Msg("Failed to parse overview response")
		return nil, err
	}

	if overview.Symbol == "" {
		log.Error().Str("symbol", symbol).Msg("Invalid symbol or no data returned")
		return nil, fmt.Errorf("invalid stock symbol")
	}

	// Check if stock exists
	stock, err := GetStockBySymbol(symbol)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		log.Error().Err(err).Str("symbol", symbol).Msg("Failed to check existing stock")
		return nil, err
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		stock = &Stock{
			Symbol: symbol,
			Name:   overview.Name,
			Sector: overview.Sector,
		}
		if _, err := stock.CreateStock(); err != nil {
			return nil, err
		}
	} else {
		stock.Name = overview.Name
		stock.Sector = overview.Sector
		if err := UpdateStock(*stock); err != nil {
			log.Error().Err(err).Str("symbol", symbol).Msg("Failed to update stock")
			return nil, err
		}
	}

	return stock, nil
}

func DeleteStock(id string) error {
	return database.DB.Where("stock_id = ?", id).Delete(&Stock{}).Error
}
