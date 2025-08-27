package controller

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/pratyush934/tradealpha/server/models"
	"github.com/pratyush934/tradealpha/server/types"
)

func GetUsersById(c echo.Context) error {

	id := c.Get("userId").(string)

	byId, err := models.GetUserById(id)

	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, types.StatusNotFound, "Not able to get the user via id", err)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Just able to do that",
		"user":    byId,
	})

}

func GetUserByEmail(c echo.Context) error {

	email := c.Get("email").(string)

	byEmail, err := models.GetUserByEmail(email)

}
