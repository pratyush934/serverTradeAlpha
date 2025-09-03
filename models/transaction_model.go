package models

import (
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/pratyush934/tradealpha/server/alphavantage"
	"github.com/pratyush934/tradealpha/server/database"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

type TransactionModel struct {
	Id          string    `gorm:"primaryKey;type:varchar(151)" json:"id"`
	UserId      string    `gorm:"not null" json:"userId"`
	PortFolioId string    `gorm:"not null" json:"portFolioId"`
	StockId     string    `gorm:"not null" json:"stockId"`
	Quantity    int       `gorm:"default:0" json:"quantity"`
	Price       float64   `gorm:"default:0" json:"price"`
	Type        string    `json:"type"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

func (t *TransactionModel) BeforeCreate(tx *gorm.DB) error {
	t.Id = uuid.New().String()
	t.CreatedAt = time.Now()
	t.UpdatedAt = time.Now()
	return nil
}

func (t *TransactionModel) BeforeUpdate(tx *gorm.DB) error {
	t.UpdatedAt = time.Now()
	return nil
}

func GetTransactionById(id string) (*TransactionModel, error) {
	var tx TransactionModel
	if err := database.DB.Where("id = ?", id).First(&tx).Error; err != nil {
		log.Error().Err(err).Msg("issue in transaction_model/GetTransactionById")
		return nil, err
	}
	return &tx, nil
}

func GetTransactionsByUserId(userId string) ([]TransactionModel, error) {
	var txs []TransactionModel
	if err := database.DB.Where("user_id = ?", userId).Find(&txs).Error; err != nil {
		log.Error().Err(err).Msg("issue in transaction_model/GetTransactionsByUserId")
		return nil, err
	}
	return txs, nil
}

func GetTransactionsByPortfolioId(portfolioId string) ([]TransactionModel, error) {
	var txs []TransactionModel
	if err := database.DB.Where("portfolio_id = ?", portfolioId).Find(&txs).Error; err != nil {
		log.Error().Err(err).Msg("issue in transaction_model/GetTransactionsByPortfolioId")
		return nil, err
	}
	return txs, nil
}

func GetTransactionsByStockId(stockId string) ([]TransactionModel, error) {
	var txs []TransactionModel
	if err := database.DB.Where("stock_id = ?", stockId).Find(&txs).Error; err != nil {
		log.Error().Err(err).Msg("issue in transaction_model/GetTransactionsByStockId")
		return nil, err
	}
	return txs, nil
}

func UpdateTransaction(tx *TransactionModel) error {
	if err := database.DB.Updates(tx).Error; err != nil {
		log.Error().Err(err).Msg("issue in transaction_model/UpdateTransaction")
		return err
	}
	return nil
}

func DeleteTransactionById(id string) error {
	if err := database.DB.Where("transaction_id = ?", id).Delete(&TransactionModel{}).Error; err != nil {
		log.Error().Err(err).Msg("issue in transaction_model/DeleteTransactionById")
		return err
	}
	return nil
}

func (t *TransactionModel) CreateTransaction(logger *zerolog.Logger) (*TransactionModel, error) {
	quote, err := alphavantage.FetchQuote(t.StockId, logger)
	if err != nil {
		log.Error().Err(err).Str("stock_id", t.StockId).Msg("Failed to fetch stock quote")
		return nil, err
	}

	price, err := strconv.ParseFloat(quote.GlobalQuote.Price, 64)
	if err != nil {
		log.Error().Err(err).Str("stock_id", t.StockId).Msg("Failed to parse stock price")
		return nil, err
	}
	t.Price = price

	if err := database.DB.Create(t).Error; err != nil {
		log.Error().Err(err).Msg("issue in transaction_model/CreateTransaction")
		return nil, err
	}
	return t, nil
}
