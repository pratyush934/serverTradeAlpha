package models

import (
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/pratyush934/tradealpha/server/alphavantage"
	"github.com/pratyush934/tradealpha/server/database"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

type PortFolio struct {
	Id              string             `gorm:"primaryKey;type:varchar(151)" json:"id"`
	UserId          string             `gorm:"not null" json:"userId"`
	Name            string             `gorm:"not null" json:"name"`
	Title           string             `gorm:"not null" json:"title"`
	TotalValue      float64            `gorm:"default:0" json:"totalValue"`
	UnRealizedGains float64            `gorm:"default:0" json:"unRealizedGains"`
	RealizedGains   float64            `gorm:"default:0" json:"realizedGains"`
	Description     string             `gorm:"not null" json:"description"`
	Transaction     []TransactionModel `gorm:"foreignKey:portFolioId" json:"transaction"`
	PortFolioStock  []PortFolioStock   `gorm:"foreignKey:portFolioId" json:"portFolioStock"`
	CreatedAt       time.Time          `json:"createdAt"`
	UpdatedAt       time.Time          `json:"updatedAt"`
}

func (p *PortFolio) BeforeCreate(tx *gorm.DB) error {
	p.Id = uuid.New().String()
	p.CreatedAt = time.Now()
	p.UpdatedAt = time.Now()

	return nil
}

func (p *PortFolio) BeforeUpdate(tx *gorm.DB) error {
	p.UpdatedAt = time.Now()
	return nil
}

func (p *PortFolio) CreatePortfolio() (*PortFolio, error) {
	if err := database.DB.Create(p).Error; err != nil {
		log.Error().Err(err).Msg("issue persist in the portfolio_model/CreatePortFolio")
		return nil, err
	}
	return p, nil
}

func GetPortFolioById(id string) (*PortFolio, error) {

	var portFolio PortFolio
	if err := database.DB.
		Preload("Transaction").
		Preload("PortFolioStock").
		Where("id = ?", id).
		First(&portFolio).Error; err != nil {
		log.Error().Err(err).Msg("issue persist in the portfolio_model/GetPortFolioById")
		return nil, err
	}
	return &portFolio, nil
}

func GetPortFolioByUserId(userId string) (*[]PortFolio, error) {

	var portfolio []PortFolio
	if err := database.DB.
		Where("user_id = ?", userId).
		Preload("Transaction").
		Preload("PortFolioStock").
		First(&portfolio).Error; err != nil {
		log.Error().Err(err).Msg("issue persist in the portfolio_model/GetPortFolioByUserId")
		return nil, err
	}
	return &portfolio, nil
}

func GetAllPortFolioStock(limit, offset int) ([]PortFolio, error) {
	var portfolio []PortFolio
	if err := database.DB.
		Preload("Transaction").
		Preload("PortFolioStock").
		Find(&portfolio).
		Error; err != nil {
		log.Error().Err(err).Msg("issue persist in the portfolio_model/GetAllPortFolioStock")
		return nil, err
	}
	return portfolio, nil
}

func UpdatePortFolio(portfolio PortFolio) error {
	return database.DB.Updates(&portfolio).Error
}

func UpdateTotalValue(id string) error {

	portfolio, err := GetPortFolioById(id)
	if err != nil {
		log.Error().Err(err).Str("portfolio_id", id).Msg("Failed to fetch portfolio")
		return err
	}

	var totalValue, unRealizedGains float64
	for _, ps := range portfolio.PortFolioStock {

		quote, err := alphavantage.FetchQuote(ps.StockId, &log.Logger)
		if err != nil {
			log.Error().Err(err).Str("stock_id", ps.StockId).Msg("Failed to fetch stock quote")
			continue
		}
		price, err := strconv.ParseFloat(quote.GlobalQuote.Price, 64)
		if err != nil {
			log.Error().Err(err).Str("stock_id", ps.StockId).Msg("Failed to parse stock price")
			continue
		}
		totalValue += float64(ps.Quantity) * price
		unRealizedGains += float64(ps.Quantity) * (price - ps.AveragePrice)
	}

	transactionsByPortfolioId, err := GetTransactionsByPortfolioId(id)

	if err != nil {
		log.Error().Err(err).Msg("there is an issue in portfolio_model/UpdateTotalValue")
		return err
	}

	var realizedGains float64
	for _, tx := range transactionsByPortfolioId {
		if tx.Type == "sell" {
			idAndPortfolioId, err := GetPortfolioStockByStockIdAndPortfolioId(tx.StockId, id)

			if err != nil || len(*idAndPortfolioId) == 0 {
				log.Error().Err(err).Msg("aha aha aha aha")
				return err
			}
			averagePrice := (*idAndPortfolioId)[0].AveragePrice
			realizedGains += float64(tx.Quantity) * (tx.Price - averagePrice)
		}
	}

	updates := map[string]interface{}{
		"total_value":      totalValue,
		"unrealized_gains": unRealizedGains,
		"realized_gains":   realizedGains,
	}

	if err := database.DB.Model(&PortFolioStock{}).Where("id = ?", id).Updates(updates).Error; err != nil {
		log.Error().Err(err).Msg("issue persist in updating the UpdateTotalValue")
		return err
	}

	log.Info().Float64("total_value", totalValue).
		Float64("unrealized_gains", unRealizedGains).
		Float64("realized_gains", realizedGains).
		Msg("You did it")

	return nil
}
