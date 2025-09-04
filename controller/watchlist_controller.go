package controller

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/pratyush934/tradealpha/server/dto"
	"github.com/pratyush934/tradealpha/server/models"
	"github.com/pratyush934/tradealpha/server/types"
	"github.com/pratyush934/tradealpha/server/util"
	"github.com/rs/zerolog/log"
)

/*
CreateWatchlistHandler: Create a new watchlist for a user.
GetWatchlistByIdHandler: Fetch a watchlist by ID, including stocks.
GetUserWatchListsHandler: Fetch all watchLists for a user.
DeleteWatchlistHandler: Delete a watchlist by ID.
AddStockToWatchlistHandler: Add a stock to a watchlist.
RemoveStockFromWatchlistHandler: Remove a stock from a watchlist.
*/

func CheckAuthorization(c echo.Context) error {

	userId := c.Get("userId").(string)

	if userId == "" {
		return util.NewAppError(http.StatusUnauthorized, types.StatusUnauthorized, "not able to get the userId", nil)
	}
	return nil
}

func CreateWatchListHandler(c echo.Context) error {

	userId := c.Get("userId").(string)

	if err := CheckAuthorization(c); err != nil {
		return err
	}

	var watchDto dto.WatchListDTO
	if err := c.Bind(&watchDto); err != nil {
		return util.NewAppError(http.StatusNotFound, types.StatusNotFound, "not able to bind the watchDto", err)
	}

	newWatchDTO := models.WatchListModel{
		UserId:         userId,
		Name:           watchDto.Name,
		Description:    watchDto.Description,
		WatchListStock: make([]models.WatchListStockModel, 0),
	}

	watchListModel, err := newWatchDTO.Create()

	if err != nil {
		log.Error().Err(err).Msg("issue persist in the watchlist_controller while create")
		return util.NewAppError(http.StatusBadRequest, types.StatusBadRequest, "not able to create the stuff", err)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"watchList": watchListModel,
	})
}

func GetWatchlistByIdHandler(c echo.Context) error {

	if err := CheckAuthorization(c); err != nil {
		return err
	}

	watchId := c.Param("watchId")

	watchListById, err := models.GetWatchListById(watchId)

	if err != nil {
		return util.NewAppError(http.StatusBadRequest, types.StatusBadRequest, "not able to get the watchListById", err)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"watchListById": watchListById,
	})
}

func DeleteWatchlistHandler(c echo.Context) error {

	if err := CheckAuthorization(c); err != nil {
		return err
	}

	watchId := c.Param("watchId")

	if watchId == "" {
		return util.NewAppError(http.StatusBadRequest, types.StatusBadRequest, "not able to get the watchId", nil)
	}

	if err := models.DeleteWatchList(watchId); err != nil {
		return util.NewAppError(http.StatusBadRequest, types.StatusBadRequest, "not able to delete the watchList", err)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Deleted successfully",
	})
}

func AddStockToWatchlistHandler(c echo.Context) error {
	if err := CheckAuthorization(c); err != nil {
		return err
	}

	watchId := c.QueryParam("watchListId")
	stockId := c.QueryParam("stockId")

	if watchId == "" || stockId == "" {
		return util.NewAppError(http.StatusBadRequest, types.StatusBadGateway, "not able to get the watchId or stockId", nil)
	}
	
	if err := models.AddStockToWatchlist(watchId, stockId); err != nil {
		return util.NewAppError(http.StatusBadRequest, types.StatusBadRequest, "not able to add the stock to watchlist", err)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Added stock to the watchList",
	})
}

func RemoveStockFromWatchlistHandler(c echo.Context) error {

	if err := CheckAuthorization(c); err != nil {
		return err
	}

	watchId := c.QueryParam("watchListId")
	stockId := c.QueryParam("stockId")

	if watchId == "" || stockId == "" {
		return util.NewAppError(http.StatusBadRequest, types.StatusBadGateway, "not able to get the watchId or stockId", nil)
	}

	if err := models.RemoveStockFromWatchlist(watchId, stockId); err != nil {
		return util.NewAppError(http.StatusBadRequest, types.StatusBadRequest, "not able to add the stock to watchlist", err)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Removed the stock from WatchList",
	})
}
