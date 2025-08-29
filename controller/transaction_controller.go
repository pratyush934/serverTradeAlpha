package controller

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/pratyush934/tradealpha/server/dto"
	"github.com/pratyush934/tradealpha/server/models"
	"github.com/pratyush934/tradealpha/server/types"
	"github.com/pratyush934/tradealpha/server/util"
)

func CreateTransaction(c echo.Context) error {

	userId := c.Get("userId").(string)

	if userId == "" {
		return util.NewAppError(http.StatusUnauthorized, types.StatusUnauthorized, "not able to get the userId", nil)
	}

	stockId := c.Param("stockId")
	portId := c.Param("portFolioId")

	if stockId == "" || portId == "" {
		return util.NewAppError(http.StatusBadRequest, types.StatusBadRequest, "not able to get the stockId or protId", nil)
	}

	var transaction dto.TransactionDTO

	if err := c.Bind(&transaction); err != nil {
		return util.NewAppError(http.StatusNotFound, types.StatusNotFound, "not able to bind the transaction", err)
	}

	newTransaction := models.TransactionModel{
		UserId:      userId,
		StockId:     stockId,
		PortFolioId: portId,
		Quantity:    transaction.Quantity,
		Price:       transaction.Price,
		Status:      transaction.Status,
		Type:        transaction.Type,
	}

	createTransaction, err := newTransaction.CreateTransaction()

	if err != nil {
		return util.NewAppError(http.StatusBadRequest, types.StatusBadRequest, "not able to create transaction", err)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"transaction": createTransaction,
	})
}

func GetPortFolioTransactionByUserId(c echo.Context) error {
	userId := c.Get("userId").(string)

	if userId == "" {
		return util.NewAppError(http.StatusUnauthorized, types.StatusUnauthorized, "not able to get the userId", nil)
	}

	transactionsByUserId, err := models.GetTransactionsByUserId(userId)

	if err != nil {
		return util.NewAppError(http.StatusBadRequest, types.StatusBadRequest, "not able to get the transaction", err)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"transactionByUserId": transactionsByUserId,
	})
}
