package controller

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/pratyush934/tradealpha/server/dto"
	"github.com/pratyush934/tradealpha/server/models"
	"github.com/pratyush934/tradealpha/server/types"
	"github.com/pratyush934/tradealpha/server/util"
)

func CreateStock(c echo.Context) error {

	userId := c.Get("userId").(string)

	if userId == "" {
		return util.NewAppError(http.StatusUnauthorized, types.StatusUnauthorized, "not able to get the userid in CreateStock", nil)
	}

	var stockDto dto.StockDTO

	if err := c.Bind(&stockDto); err != nil {
		return util.NewAppError(http.StatusNotFound, types.StatusNotFound, "not able to bind the stockDto", err)
	}

	newStock := models.Stock{
		Name:           stockDto.Name,
		Sector:         stockDto.Sector,
		Price:          stockDto.Price,
		PortFolioStock: make([]models.PortFolioStock, 0),
		Transaction:    make([]models.TransactionModel, 0),
	}

	stock, err := newStock.CreateStock()

	if err != nil {
		return util.NewAppError(http.StatusBadRequest, types.StatusBadRequest, "not able to create the stock", err)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"stock": stock,
	})
}
