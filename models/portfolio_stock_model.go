package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/pratyush934/tradealpha/server/database"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

type PortFolioStock struct {
	Id           string    `gorm:"primaryKey; type:varchar(151)" json:"id"`
	StockId      string    `gorm:"not null" json:"stockId"`
	PortFolioId  string    `gorm:"not null" json:"portFolioId"`
	Quantity     int       `gorm:"default:0" json:"quantity"`
	AveragePrice float64   `gorm:"default:0" json:"averagePrice"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

func (p *PortFolioStock) BeforeCreate(tx *gorm.DB) error {
	p.Id = uuid.New().String()
	p.CreatedAt = time.Now()
	p.UpdatedAt = time.Now()

	return nil
}

func (p *PortFolioStock) BeforeUpdate(tx *gorm.DB) error {
	p.UpdatedAt = time.Now()
	return nil
}

func (p *PortFolioStock) CreatePortFolioStock() (*PortFolioStock, error) {
	if err := database.DB.Create(p).Error; err != nil {
		log.Error().Err(err).Msg("issue persist in portfolio_stock_model/CreatePortFolioStock")
		return nil, err
	}
	return p, nil
}

func GetPortfolioStockByStockId(stockId string) (*[]PortFolioStock, error) {
	var portfolioStock []PortFolioStock
	if err := database.DB.Where("stock_id = ?", stockId).Find(&portfolioStock).Error; err != nil {
		log.Error().Err(err).Msg("issue persist in portfolio_stock_model/GetPortFolioStockById")
		return nil, err
	}
	return &portfolioStock, nil
}

func GetPortfolioPortfolioId(pid string) (*[]PortFolioStock, error) {
	var portfolioStock []PortFolioStock
	if err := database.DB.Where("portfolio_id = ?", pid).Find(&portfolioStock).Error; err != nil {
		log.Error().Err(err).Msg("issue persist in portfolio_stock_model/GetPortFolioPortFolioId")
		return nil, err
	}
	return &portfolioStock, nil
}

func GetPortfolioStockByStockIdAndPortfolioId(sid, pid string) (*[]PortFolioStock, error) {
	var portfolioStock []PortFolioStock
	if err := database.DB.Where("portfolio_id = ? AND stock_id = ?", pid, sid).Find(&portfolioStock).Error; err != nil {
		log.Error().Err(err).Msg("issue persist in portfolio_stock_model/GetPortfolioStockByStockIdAndPortFolioId")
		return nil, err
	}
	return &portfolioStock, nil
}

func GetPortfolioStockById(id string) (*PortFolioStock, error) {
	var portfolio *PortFolioStock
	if err := database.DB.Where("portfolio_id = ?", id).Find(portfolio).Error; err != nil {
		log.Error().Err(err).Msg("issue persist in portfolio_stock_model/GetPortfolioStockById")
		return nil, err
	}
	return portfolio, nil
}

func UpdatePortfolioStock(portfolio *PortFolioStock) error {
	return database.DB.Updates(portfolio).Error
}

func DeletePortfolioById(id string) error {
	return database.DB.Where("portfolio_id = ?", id).Delete(&PortFolioStock{}).Error
}

func UpdatePortfolioStockAveragePrice(id, value string) error {
	return database.DB.Model(&PortFolioStock{}).Where("portfolio_id = ?", id).Update("average_price = ? ", value).Error
}
