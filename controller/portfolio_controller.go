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

func CreatePortfolio(c echo.Context) error {
	userId := c.Get("userId").(string)

	if userId == "" {
		return util.NewAppError(http.StatusUnauthorized, types.StatusUnauthorized, "not able to get the userId", nil)
	}

	var portfolio dto.PortFolioDTO

	if err := c.Bind(&portfolio); err != nil {
		return util.NewAppError(http.StatusBadRequest, types.StatusBadRequest, "not able to bind the portfolioId", nil)
	}

	newPortFolio := models.PortFolio{
		UserId:         userId,
		Name:           portfolio.Name,
		Title:          portfolio.Title,
		Description:    portfolio.Description,
		Transaction:    make([]models.TransactionModel, 0),
		PortFolioStock: make([]models.PortFolioStock, 0),
	}

	createPortfolio, err := newPortFolio.CreatePortfolio()

	if err != nil {
		return util.NewAppError(http.StatusBadRequest, types.StatusBadRequest, "not able to create portfolio", err)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"portfolio": createPortfolio,
	})
}

func GetPortFolioById(c echo.Context) error {
	userId := c.Get("userId").(string)

	if userId == "" {
		return util.NewAppError(http.StatusUnauthorized, types.StatusUnauthorized, "not able to get the userId , GetPortFolioById", nil)
	}

	portId := c.Param("id")

	if portId == "" {
		return util.NewAppError(http.StatusBadRequest, types.StatusBadRequest, "please provide portfolioId", nil)
	}

	portFolioById, err := models.GetPortFolioById(portId)

	if err != nil {
		return util.NewAppError(http.StatusBadRequest, types.StatusBadRequest, "not able te get the protFolio by Id", err)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"portfolio": portFolioById,
	})
}

func UpdatePortFolioStock(c echo.Context) error {

	userId := c.Get("userId").(string)

	if userId == "" {
		return util.NewAppError(http.StatusUnauthorized, types.StatusUnauthorized, "not able to get the userid, UpdateThePortFolio", nil)
	}

	portId := c.Param("id")

	if portId == "" {
		return util.NewAppError(http.StatusBadRequest, types.StatusBadRequest, "not able to get the portId, UpdatePortfolio", nil)
	}

	var portfolio dto.PortFolioDTO

	if err := c.Bind(&portfolio); err != nil {
		return util.NewAppError(http.StatusBadRequest, types.StatusBadRequest, "not able to bind the portfolio, UpdatePortFolio", err)
	}

	portFolioById, err := models.GetPortFolioById(portId)

	if err != nil || portFolioById == nil {
		return util.NewAppError(http.StatusBadRequest, types.StatusBadRequest, "not able te get the protFolio by Id", err)
	}

	portFolioById.Title = portfolio.Title
	portFolioById.Name = portfolio.Name
	portFolioById.Description = portfolio.Description

	if err := models.UpdatePortFolio(*portFolioById); err != nil {
		return util.NewAppError(http.StatusBadRequest, types.StatusBadRequest, "not able to get the Update portfolio", err)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"update_status": http.StatusOK,
	})

}

func DeletePortFolio(c echo.Context) error {

	userId := c.Get("userId").(string)

	if userId == "" {
		return util.NewAppError(http.StatusUnauthorized, types.StatusUnauthorized, "not able to get the userId", nil)
	}

	portId := c.Param("id")

	if err := models.DeletePortfolioById(portId); err != nil {
		return util.NewAppError(http.StatusBadRequest, types.StatusBadRequest, "not able to delete the portfolio by portId", err)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "PortFolioDelete successfully",
	})

}

func UpdatePortFolioTotalValue(c echo.Context) error {

	userId := c.Get("userId").(string)

	if userId == "" {
		return util.NewAppError(http.StatusUnauthorized, types.StatusUnauthorized, "not able to get the userId", nil)
	}

	v := c.Param("value")

	value, _ := strconv.Atoi(v)

	if err := models.UpdateTotalValue(userId, value); err != nil {
		return util.NewAppError(http.StatusBadRequest, types.StatusBadRequest, "not able to update_value in UpdatePortFolioTotalValue", err)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"value": value,
	})

}
