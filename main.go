package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/pratyush934/tradealpha/server/controller"
	"github.com/pratyush934/tradealpha/server/types"
	"github.com/pratyush934/tradealpha/server/util"
	"github.com/rs/zerolog"
)

func LoadDb() {

}

func Server() {

	e := echo.New()

	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()
	e.Use(util.ErrorHandleMiddleWare(&logger))

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "hello")
	})

	e.GET("/test", func(c echo.Context) error {
		return util.NewAppError(http.StatusOK, types.StatusOK, "It is working bro", fmt.Errorf("first time , this is first time"))
	})

	e.POST("/login", controller.LoginController)

	e.GET("/api/stocks/search", util.SearchStockHandler(&logger))
	e.GET("/api/stocks/:symbol/quote", util.GetStockQuoteHandler(&logger))
	e.GET("/api/stocks/:symbol/intraday", util.GetIntradayDataHandler(&logger))

	_ = e.Start(":8080")

}

func Config() {

}

func main() {
	Server()
}
