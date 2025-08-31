package controller

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/pratyush934/tradealpha/server/dto"
	"github.com/pratyush934/tradealpha/server/models"
	"github.com/pratyush934/tradealpha/server/types"
	"github.com/pratyush934/tradealpha/server/util"
	"github.com/rs/zerolog/log"
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

	candidate, err := models.GetUserByEmail(email)

	if err == nil {

		err2 := models.UpdateLastLogin(email, time.Now())

		if err2 != nil {
			log.Error().Err(err).Msg("Please check the UpdateLastLogin")
			return util.NewAppError(http.StatusBadRequest, types.StatusBadRequest, "Please check the err2", err2)
		}

		token, err3 := util.CreateToken(candidate)

		if err3 != nil {
			log.Error().Err(err3).Msg("not able to generate token , err3")
			return util.NewAppError(http.StatusBadRequest, types.StatusBadRequest, "Please check the err3", err3)
		}

		return c.JSON(http.StatusOK, map[string]interface{}{
			"message": "the user already exist, Login Successful",
			"email":   token,
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
