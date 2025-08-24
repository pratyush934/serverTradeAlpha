package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/pratyush934/tradealpha/server/types"
	"github.com/pratyush934/tradealpha/server/utils"
	"github.com/rs/zerolog"
)

func LoadDb() {

}

func Server() {

	e := echo.New()

	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()
	e.Use(utils.ErrorHandleMiddleWare(&logger))

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "hello")
	})

	e.GET("/test", func(c echo.Context) error {
		return utils.NewAppError(http.StatusOK, types.StatusOK, "It is working bro", fmt.Errorf("first time , this is first time"))
	})

	_ = e.Start(":8080")

}

func Config() {

}

func main() {
	Server()
}
