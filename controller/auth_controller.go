package controller

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/pratyush934/tradealpha/server/dto"
	"github.com/pratyush934/tradealpha/server/models"
	"github.com/pratyush934/tradealpha/server/util"
)

func LoginController(c echo.Context) error {
	var login dto.LoginModel

	if err := c.Bind(&login); err != nil {
		return c.JSON(http.StatusNoContent, map[string]interface{}{
			"error":   err,
			"message": "Look at the LoginController",
		})
	}

	email := login.Email

	_, err := models.GetUserByEmail(email)
	if err == nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "the user already exist",
			"email":   email,
		})
	}

	newUser := models.User{
		Name:         login.Name,
		Email:        login.Email,
		ProfileImage: login.Image,
		OAuthId:      login.OAuthId,
		Provider:     login.Provider,
	}

	user, err := newUser.CreateUser()

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"message": "Not able to create the user",
			"Error":   err,
		})
	}

	token, err := util.CreateToken(user)

	if err != nil {

		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"message": "Not able to create the token",
			"Error":   err,
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"user":  user,
		"token": token,
	})
}
