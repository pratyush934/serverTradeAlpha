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

	stockId := c.QueryParam("stockId")
	portId := c.QueryParam("portFolioId")

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

func GetTransactionByUserId(c echo.Context) error {
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

func GetPortFolioTransactionById(c echo.Context) error {
	userId := c.Get("userId").(string)

	if userId == "" {
		return util.NewAppError(http.StatusUnauthorized, types.StatusUnauthorized, "not able to get the userId", nil)
	}

	transactionId := c.Param("transId")

	if transactionId == "" {
		return util.NewAppError(http.StatusBadRequest, types.StatusBadRequest, "transactionId is empty", nil)
	}

	transactionById, err := models.GetTransactionById(transactionId)

	if err != nil {
		return util.NewAppError(http.StatusBadRequest, types.StatusBadRequest, "not able to get the transactionId", err)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"transactionById": transactionById,
	})
}

func GetPortFolioTransactionByStockId(c echo.Context) error {
	userId := c.Get("userId").(string)

	if userId == "" {
		return util.NewAppError(http.StatusUnauthorized, types.StatusUnauthorized, "not able to get the userId", nil)
	}

	stockId := c.Param("stockId")

	if stockId == "" {
		return util.NewAppError(http.StatusBadRequest, types.StatusBadRequest, "transactionId is empty", nil)
	}

	transactionById, err := models.GetTransactionsByStockId(stockId)

	if err != nil {
		return util.NewAppError(http.StatusBadRequest, types.StatusBadRequest, "not able to get the transactionId", err)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"transactionById": transactionById,
	})
}

func GetPortFolioTransactionByPortFolioId(c echo.Context) error {
	userId := c.Get("userId").(string)

	if userId == "" {
		return util.NewAppError(http.StatusUnauthorized, types.StatusUnauthorized, "not able to get the userId", nil)
	}

	portId := c.Param("portId")

	if portId == "" {
		return util.NewAppError(http.StatusBadRequest, types.StatusBadRequest, "transactionId is empty", nil)
	}

	transactionById, err := models.GetTransactionsByPortfolioId(portId)

	if err != nil {
		return util.NewAppError(http.StatusBadRequest, types.StatusBadRequest, "not able to get the transactionId", err)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"transactionById": transactionById,
	})
}

func DeleteTransactionByUserId(c echo.Context) error {
	userId := c.Get("userId").(string)

	if userId == "" {
		return util.NewAppError(http.StatusUnauthorized, types.StatusUnauthorized, "not able to get the userId", nil)
	}

	transactionId := c.Param("transId")

	if transactionId == "" {
		return util.NewAppError(http.StatusBadRequest, types.StatusBadRequest, "transactionId is empty", nil)
	}

	if err := models.DeleteTransactionById(transactionId); err != nil {
		return util.NewAppError(http.StatusBadRequest, types.StatusBadRequest, "not able to delete the transactions", err)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "DeletedSuccessfully",
	})
}

func UpdateTransaction(c echo.Context) error {
	userId := c.Get("userId").(string)

	if userId == "" {
		return util.NewAppError(http.StatusUnauthorized, types.StatusUnauthorized, "not able to get the userId", nil)
	}

	transactionId := c.Param("transId")

	if transactionId == "" {
		return util.NewAppError(http.StatusBadRequest, types.StatusBadRequest, "transactionId is empty", nil)
	}

	var transaction dto.TransactionDTO

	if err := c.Bind(&transaction); err != nil {
		return util.NewAppError(http.StatusNotFound, types.StatusNotFound, "not able to bind the transaction", err)
	}

	transactionById, err := models.GetTransactionById(transactionId)

	if err != nil {
		return util.NewAppError(http.StatusBadRequest, types.StatusBadRequest, "not able to get the transactionById", err)
	}

	if transactionById.UserId != userId {
		return util.NewAppError(http.StatusBadRequest, types.StatusBadRequest, "not the correct user", nil)
	}

	transactionById.Type = transaction.Type
	transactionById.Status = transaction.Status
	transactionById.Quantity = transaction.Quantity
	transactionById.Price = transaction.Price

	if err := models.UpdateTransaction(transactionById); err != nil {
		return util.NewAppError(http.StatusBadRequest, types.StatusBadRequest, "not able to update the transaction", err)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "updated transaction",
	})
}

func GetTransactionsByStockId(c echo.Context) error {

	userId := c.Get("userId").(string)

	if userId == "" {
		return util.NewAppError(http.StatusUnauthorized, types.StatusUnauthorized, "not able to get the userId", nil)
	}

	stockId := c.Param("stockId")

	byStockId, err := models.GetTransactionsByStockId(stockId)

	if err != nil {
		return util.NewAppError(http.StatusOK, types.StatusBadRequest, "not able to get transaction by stockId", err)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"stockById": byStockId,
	})

}
