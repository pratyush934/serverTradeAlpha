package jwtpackage

import (
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/pratyush934/tradealpha/server/types"
	"github.com/pratyush934/tradealpha/server/util"
	"github.com/rs/zerolog/log"
)

func ValidateUserMiddleWare() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {

			token, err := GetToken(c)
			if err != nil {
				log.Error().Err(err).Msg("There is an issue in the ValidateUserMiddleware")
				return util.NewAppError(http.StatusBadRequest, types.StatusBadRequest, "Not able to GetToken in the middleware", err)
			}

			if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
				c.Set("userId", claims["id"])
				c.Set("email", claims["email"])
				c.Set("name", claims["name"])
				c.Set("role", claims["role"])
			} else {
				return util.NewAppError(http.StatusInternalServerError, types.StatusInternalServerError, "Not able to set claims in context", nil)
			}

			return next(c)
		}
	}
}

func ValidateAdminMiddleWare() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {

			token, err := GetToken(c)
			if err != nil {
				log.Error().Err(err).Msg("there is an issue in the ValidateAdminMiddleware")
				return util.NewAppError(http.StatusBadRequest, types.StatusBadRequest, "Not able to GetToken in the middleware", err)
			}

			claims, ok := token.Claims.(jwt.MapClaims)

			role := claims["role"].(int64)

			if !ok || !token.Valid {
				return util.NewAppError(http.StatusNotFound, types.StatusNotFound, "Token is not valid", nil)
			}

			if role != 2 {
				return util.NewAppError(http.StatusUnauthorized, types.StatusUnauthorized, "Not authorized as the user is not admin", nil)
			}

			c.Set("userId", claims["id"])
			c.Set("role", claims["role"])
			c.Set("name", claims["name"])
			c.Set("email", claims["email"])

			return next(c)
		}
	}
}
