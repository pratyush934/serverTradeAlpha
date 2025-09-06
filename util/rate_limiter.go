package util

import (
	"net/http"
	"sync"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/pratyush934/tradealpha/server/types"
	"github.com/rs/zerolog/log"
)

var requestLog = make(map[string][]time.Time)
var mu sync.Mutex

const (
	MaxLimit   = 5
	TimeWindow = time.Minute
)

func RateLimiter() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {

			ip := c.Request().RemoteAddr

			if ip == "" {
				log.Error().Str("there is not ip address we got here : ", ip).Msg("Please provide the remote addr in RateLimiter")
				return NewAppError(http.StatusBadRequest, types.StatusBadRequest, "look at the rate-limiter", nil)
			}

			mu.Lock()

			requests := requestLog[ip]
			currentTime := time.Now()
			var validateTime []time.Time

			for _, t := range requests {
				if currentTime.Sub(t) < TimeWindow {
					validateTime = append(validateTime, t)
				}
			}

			requestLog[ip] = append(validateTime, currentTime)
			mu.Unlock()

			if len(validateTime) > MaxLimit {
				return NewAppError(http.StatusBadGateway, types.StatusBadGateway, "you have exceeded the rate limiter", nil)
			}

			return next(c)
		}
	}
}
