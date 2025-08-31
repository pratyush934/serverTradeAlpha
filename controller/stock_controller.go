package controller

import (
	"net/http"
	"strconv"

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

func GetStockById(c echo.Context) error {

	userId := c.Get("userId").(string)

	if userId == "" {
		return util.NewAppError(http.StatusUnauthorized, types.StatusUnauthorized, "not able to get the userId", nil)
	}

	stockId := c.Param("id")

	stockById, err := models.GetStockById(stockId)

	if err != nil {
		return util.NewAppError(http.StatusBadRequest, types.StatusBadRequest, "not able to get the stockById", err)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"stock-by-id": stockById,
	})
}

func GetStockBySector(c echo.Context) error {

	userId := c.Get("userId").(string)

	if userId == "" {
		return util.NewAppError(http.StatusUnauthorized, types.StatusUnauthorized, "not able to get the userId", nil)
	}

	sector := c.QueryParam("sector")
	limitStr := c.QueryParam("limit")
	offSetStr := c.QueryParam("offSet")

	limit, _ := strconv.Atoi(limitStr)
	offSet, _ := strconv.Atoi(offSetStr)

	stockBySector, err := models.GetStockBySector(sector, limit, offSet)

	if err != nil {
		return util.NewAppError(http.StatusBadRequest, types.StatusBadRequest, "not able to get the stockBySector", err)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"stock-by-sector": stockBySector,
	})
}

func GetAllStocks(c echo.Context) error {

	userId := c.Get("userId").(string)

	if userId == "" {
		return util.NewAppError(http.StatusUnauthorized, types.StatusUnauthorized, "not able to get the userId", nil)
	}
	limitStr := c.QueryParam("limit")
	offSetStr := c.QueryParam("offSet")

	limit, _ := strconv.Atoi(limitStr)
	offSet, _ := strconv.Atoi(offSetStr)

	getAllStocks, err := models.GetAllStocks(limit, offSet)

	if err != nil {
		return util.NewAppError(http.StatusBadRequest, types.StatusBadRequest, "not able to get all the stocks", err)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"stocks": getAllStocks,
	})

}

func UpdateStock(c echo.Context) error {
	userId := c.Get("userId").(string)

	if userId == "" {
		return util.NewAppError(http.StatusUnauthorized, types.StatusUnauthorized, "not able to get the userId", nil)
	}

	stockId := c.Param("stockId")

	var stockDTO dto.StockDTO

	if err := c.Bind(&stockDTO); err != nil {
		return util.NewAppError(http.StatusBadRequest, types.StatusBadRequest, "not able to bind the stockDTO", nil)
	}

	stockById, err := models.GetStockById(stockId)

	if err != nil {
		return util.NewAppError(http.StatusBadRequest, types.StatusBadRequest, "not able to get the stockById", err)
	}
	stockById.Name = stockDTO.Name
	stockById.Price = stockDTO.Price
	stockById.Sector = stockDTO.Sector

	if err := models.UpdateStock(*stockById); err != nil {
		return util.NewAppError(http.StatusBadRequest, types.StatusBadRequest, "not able to get the UpdateStock", err)
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "updated the stock",
	})
}

func DeleteStock(c echo.Context) error {
	userId := c.Get("userId").(string)

	if userId == "" {
		return util.NewAppError(http.StatusUnauthorized, types.StatusUnauthorized, "not able to get the userId", nil)
	}

	stockId := c.Param("stockId")

	if err := models.DeleteStock(stockId); err != nil {
		return util.NewAppError(http.StatusBadRequest, types.StatusBadRequest, "not able to delete the stock", err)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Deleted OK",
	})

}
