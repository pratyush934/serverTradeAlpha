package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/pratyush934/tradealpha/server/database"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

type Stock struct {
	Id             string             `gorm:"primaryKey;type:varchar(151)" json:"id"`
	Name           string             `gorm:"not null" json:"name"`
	Sector         string             `json:"sector"`
	Price          float64            `gorm:"not null" json:"price"`
	PortFolioStock []PortFolioStock   `gorm:"foreignKey:stockId" json:"portFolioStock"`
	Transaction    []TransactionModel `gorm:"foreignKey:stockId" json:"transaction"`
	CreatedAt      time.Time          `json:"createdAt"`
	UpdatedAt      time.Time          `json:"updatedAt"`
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
		Preload("PortFolioStock").
		Preload("Transaction").
		First(&stock).Error; err != nil {
		log.Error().Err(err).Msg("issue persist in stock_model/GetStockById")
		return nil, err
	}
	return &stock, nil
}

func GetStockBySector(sector string) (*[]Stock, error) {
	var stock []Stock
	if err := database.DB.Where("sector = ?", sector).
		Preload("PortFolioStock").
		Preload("Transaction").
		Find(&stock).Error; err != nil {
		log.Error().Err(err).Msg("issue persist in stock_model/GetStockBySector")
		return nil, err
	}
	return &stock, nil
}

func GetAllStocks(limit, offset int) ([]Stock, error) {
	var stock []Stock
	if err := database.DB.
		Preload("PortFolioStock").
		Preload("Transaction").
		Limit(limit).Offset(offset).
		Find(&stock).Error; err != nil {
		log.Error().Err(err).Msg("issue persist in stock_model/GetAllStocks")
		return nil, err
	}
	return stock, nil
}

func UpdateStock(stock Stock) error {
	return database.DB.Updates(&stock).Error
}

func DeleteStock(id string) error {
	return database.DB.Where("stock_id = ?", id).Delete(&Stock{}).Error
}
