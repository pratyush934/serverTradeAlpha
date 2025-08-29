package controller

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/pratyush934/tradealpha/server/dto"
	"github.com/pratyush934/tradealpha/server/models"
	"github.com/pratyush934/tradealpha/server/types"
	"github.com/pratyush934/tradealpha/server/util"
)

func AddPortfolioStock(c echo.Context) error {

	userId := c.Get("userId")

	if userId == "" {
		return util.NewAppError(http.StatusUnauthorized, types.StatusUnauthorized, "not able to get the userId in AddPortfolioStock", nil)
	}

	stockIdStr := c.Param("stockId")
	portIdStr := c.Param("portId")

	var portFolioStockDTO dto.PortFolioStockDTO

	if err := c.Bind(&portFolioStockDTO); err != nil {
		return util.NewAppError(http.StatusBadRequest, types.StatusBadRequest, "not able to bind the portfoliodto", err)
	}

	newPortFolioStock := models.PortFolioStock{
		StockId:      stockIdStr,
		PortFolioId:  portIdStr,
		Quantity:     portFolioStockDTO.Quantity,
		AveragePrice: portFolioStockDTO.AveragePrice,
	}

	stock, err := newPortFolioStock.CreatePortFolioStock()

	if err != nil {
		return util.NewAppError(http.StatusBadRequest, types.StatusBadRequest, " not able to create portfolio stock", err)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"portfolio stock": stock,
	})
}

func UpdatePortfolioStock(c echo.Context) error {
	userId := c.Get("userId")

	if userId == "" {
		return util.NewAppError(http.StatusUnauthorized, types.StatusUnauthorized, "not able to get the userId in AddPortfolioStock", nil)
	}

	portfolioStockIdStr := c.Param("id")

	var portFolioStockDTO dto.PortFolioStockDTO

	if err := c.Bind(&portFolioStockDTO); err != nil {
		return util.NewAppError(http.StatusBadRequest, types.StatusBadRequest, "not able to bind the portfolio", err)
	}

	portfolioPortfolioId, err := models.GetPortfolioStockById(portfolioStockIdStr)

	if err != nil || portfolioPortfolioId == nil {
		return util.NewAppError(http.StatusBadRequest, types.StatusBadRequest, " not able to create portfolio stock", err)
	}

	portfolioPortfolioId.Quantity = portFolioStockDTO.Quantity
	portfolioPortfolioId.AveragePrice = portFolioStockDTO.AveragePrice

	err = models.UpdatePortfolioStock(portfolioPortfolioId)

	if err != nil {
		return util.NewAppError(http.StatusBadRequest, types.StatusBadRequest, " not able to create portfolio stock", err)
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"portfolio stock": "Updated successfully",
	})
}

func RemovePortfolioStock(c echo.Context) error {

	userId := c.Get("userId").(string)

	if userId == "" {
		return util.NewAppError(http.StatusUnauthorized, types.StatusUnauthorized, "not able to get the userId", nil)
	}

	portStockId := c.Param("id")

	if portStockId == "" {
		return util.NewAppError(http.StatusBadRequest, types.StatusBadRequest, "not able to get the portStockId", nil)
	}

	if err := models.DeletePortfolioById(portStockId); err != nil {
		return util.NewAppError(http.StatusBadRequest, types.StatusBadRequest, "not able to delete portfolio by its id", err)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "delete successfully",
	})
}

func GetPortFolioStocks(c echo.Context) error {

	userId := c.Get("userId").(string)

	if userId == "" {
		return util.NewAppError(http.StatusUnauthorized, types.StatusUnauthorized, "not able to get the userId", nil)
	}

	portStockId := c.Param("id")

	if portStockId == "" {
		return util.NewAppError(http.StatusBadRequest, types.StatusBadRequest, "not able to get the portStockId", nil)
	}

	portfolioPortfolioId, err := models.GetPortfolioPortfolioId(portStockId)

	if err != nil {
		return util.NewAppError(http.StatusBadRequest, types.StatusBadRequest, "not able to get all the portfolios-stock", err)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"portfolio-stock-by-id": portfolioPortfolioId,
	})
}
