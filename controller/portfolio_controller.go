package controller

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/pratyush934/tradealpha/server/dto"
	"github.com/pratyush934/tradealpha/server/models"
	"github.com/pratyush934/tradealpha/server/types"
	"github.com/pratyush934/tradealpha/server/util"
)

/*
CreatePortfolio

GetPortfolioById

GetUserPortfolios

UpdatePortfolio

DeletePortfolio
*/

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
	return nil
}
