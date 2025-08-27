package util

import (
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/pratyush934/tradealpha/server/types"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

type AppError struct {
	Status        int    `json:"status"`
	Code          string `json:"code"`
	Message       string `json:"message"`
	InternalError error  `json:"-"`
}

func (a *AppError) Error() string {
	if a.InternalError != nil {
		return a.Message + " : " + a.InternalError.Error()
	}
	return a.Message
}

func NewAppError(status int, code, message string, internalError error) *AppError {
	return &AppError{
		Status:        status,
		Code:          code,
		Message:       message,
		InternalError: internalError,
	}
}

type ErrorResponse struct {
	Status  int    `json:"status"`
	Code    string `json:"code"`
	Message string `json:"message"`
}

func ErrorHandleMiddleWare(logger *zerolog.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			reqId := c.Request().Header.Get(echo.HeaderXRequestID)
			if reqId == "" {
				reqId = uuid.New().String()
				c.Request().Header.Set(echo.HeaderXRequestID, reqId)
			}

			log := logger.With().
				Str("request_id", reqId).
				Str("method", c.Request().Method).
				Str("path", c.Request().URL.Path).
				Logger()

			c.Set("logger", &log)

			defer func() {
				if r := recover(); r != nil {
					log.
						Error().
						Interface("panic", r).
						Msg("there is an unexpected panic, please look at the issue")

					_ = c.JSON(http.StatusInternalServerError, ErrorResponse{
						Status:  http.StatusInternalServerError,
						Code:    types.StatusInternalServerError,
						Message: "There is an issue in recover in the ErrorHandleMiddleware, please look at this!!",
					})
				}
			}()

			err := next(c)
			if err == nil {
				return nil
			}

			log = *c.Get("logger").(*zerolog.Logger)
			var appError *AppError

			if errors.As(err, &appError) {
				logEvent := log.Info()
				if appError.Status >= http.StatusInternalServerError {
					logEvent = log.Error()
				}

				logEvent.
					Int("status", appError.Status).
					Str("code", appError.Code).
					Str("message", appError.Message).
					Err(appError.InternalError).
					Msg("Handle Application Message Error")

				return c.JSON(appError.Status, ErrorResponse{
					Message: appError.Message,
					Status:  appError.Status,
					Code:    appError.Code,
				})

			}

			if errors.As(err, &gorm.ErrRecordNotFound) {
				log.Warn().
					Err(err).
					Msg("not able to find the record")

				return c.JSON(http.StatusNotFound, ErrorResponse{
					Message: "Not able to found stuff via gorm",
					Status:  http.StatusNotFound,
					Code:    types.StatusNotFound,
				})
			}

			log.Error().
				Err(err).
				Msg("Unexpected Error")

			return c.JSON(http.StatusInternalServerError, ErrorResponse{
				Status:  http.StatusInternalServerError,
				Code:    types.StatusInternalServerError,
				Message: "Please look at the ErrorHandlingMiddleWare Part1",
			})

		}
	}
}
