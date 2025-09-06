package util

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func CORSHandler() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {

			c.Request().Header.Set("Access-Control-Allow-Origin", "*")
			c.Request().Header.Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
			c.Request().Header.Set("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization, X-Requested-With")

			if c.Request().Method != http.MethodOptions {
				return c.String(http.StatusOK, "Look at the CORSHandler")
			}

			return next(c)
		}
	}
}
