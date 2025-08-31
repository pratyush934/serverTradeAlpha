package util

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/pratyush934/tradealpha/server/types"
	"github.com/rs/zerolog"
)

const (
	AlphaVantageBaseURL = "https://www.alphavantage.co/query"
)

// QuoteResponse represents the GLOBAL_QUOTE API response
type QuoteResponse struct {
	GlobalQuote struct {
		Symbol    string `json:"01. symbol"`
		Price     string `json:"05. price"`
		Volume    string `json:"06. volume"`
		Timestamp string `json:"07. latest trading day"`
	} `json:"Global Quote"`
}

// IntradayResponse represents the TIME_SERIES_INTRADAY API response
type IntradayResponse struct {
	MetaData struct {
		Symbol   string `json:"2. Symbol"`
		Interval string `json:"4. Interval"`
	} `json:"Meta Data"`
	TimeSeries map[string]struct {
		Open   string `json:"1. open"`
		High   string `json:"2. high"`
		Low    string `json:"3. low"`
		Close  string `json:"4. close"`
		Volume string `json:"5. volume"`
	} `json:"Time Series (1min)"` // Adjust for other intervals if needed
}

// SearchResponse represents the SYMBOL_SEARCH API response
type SearchResponse struct {
	BestMatches []struct {
		Symbol      string `json:"1. symbol"`
		Name        string `json:"2. name"`
		Type        string `json:"3. type"`
		Region      string `json:"4. region"`
		MarketOpen  string `json:"5. marketOpen"`
		MarketClose string `json:"6. marketClose"`
		Timezone    string `json:"7. timezone"`
		Currency    string `json:"8. currency"`
		MatchScore  string `json:"9. matchScore"`
	} `json:"bestMatches"`
}

// FetchQuote retrieves the current stock quote for a symbol
func FetchQuote(symbol string, logger *zerolog.Logger) (*QuoteResponse, error) {
	url := fmt.Sprintf("%s?function=GLOBAL_QUOTE&symbol=%s&apikey=%s", AlphaVantageBaseURL, symbol, os.Getenv("ALPHA_VANTAGE_KEY"))
	resp, err := http.Get(url)
	if err != nil {
		logger.Error().Err(err).Str("symbol", symbol).Msg("Failed to fetch quote from Alpha Vantage")
		return nil, NewAppError(http.StatusInternalServerError, types.StatusInternalServerError, "Failed to fetch stock quote", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		logger.Error().Int("status", resp.StatusCode).Str("symbol", symbol).Msg("Alpha Vantage API returned non-200 status")
		return nil, NewAppError(http.StatusInternalServerError, types.StatusInternalServerError, "Alpha Vantage API error", nil)
	}

	var quote QuoteResponse
	if err := json.NewDecoder(resp.Body).Decode(&quote); err != nil {
		logger.Error().Err(err).Str("symbol", symbol).Msg("Failed to parse quote response")
		return nil, NewAppError(http.StatusInternalServerError, types.StatusInternalServerError, "Failed to parse stock quote", err)
	}

	if quote.GlobalQuote.Symbol == "" {
		logger.Error().Str("symbol", symbol).Msg("Invalid symbol or no data returned")
		return nil, NewAppError(http.StatusBadRequest, types.StatusBadRequest, "Invalid stock symbol", nil)
	}

	return &quote, nil
}

// FetchIntraday retrieves intraday time series data for a symbol
func FetchIntraday(symbol, interval string, logger *zerolog.Logger) (*IntradayResponse, error) {
	if !isValidInterval(interval) {
		logger.Error().Str("interval", interval).Msg("Invalid interval for intraday data")
		return nil, NewAppError(http.StatusBadRequest, types.StatusBadRequest, "Invalid interval", nil)
	}

	url := fmt.Sprintf("%s?function=TIME_SERIES_INTRADAY&symbol=%s&interval=%s&apikey=%s&extended_hours=true", AlphaVantageBaseURL, symbol, interval, os.Getenv("ALPHA_VANTAGE_KEY"))
	resp, err := http.Get(url)
	if err != nil {
		logger.Error().Err(err).Str("symbol", symbol).Msg("Failed to fetch intraday data from Alpha Vantage")
		return nil, NewAppError(http.StatusInternalServerError, types.StatusInternalServerError, "Failed to fetch intraday data", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		logger.Error().Int("status", resp.StatusCode).Str("symbol", symbol).Msg("Alpha Vantage API returned non-200 status")
		return nil, NewAppError(http.StatusInternalServerError, types.StatusInternalServerError, "Alpha Vantage API error", nil)
	}

	var intraday IntradayResponse
	if err := json.NewDecoder(resp.Body).Decode(&intraday); err != nil {
		logger.Error().Err(err).Str("symbol", symbol).Msg("Failed to parse intraday response")
		return nil, NewAppError(http.StatusInternalServerError, types.StatusInternalServerError, "Failed to parse intraday data", err)
	}

	if intraday.MetaData.Symbol == "" {
		logger.Error().Str("symbol", symbol).Msg("Invalid symbol or no data returned")
		return nil, NewAppError(http.StatusBadRequest, types.StatusBadRequest, "Invalid stock symbol", nil)
	}

	return &intraday, nil
}

// SearchSymbol searches for stocks by keyword
func SearchSymbol(keyword string, logger *zerolog.Logger) (*SearchResponse, error) {
	url := fmt.Sprintf("%s?function=SYMBOL_SEARCH&keywords=%s&apikey=%s", AlphaVantageBaseURL, keyword, os.Getenv("ALPHA_VANTAGE_KEY"))
	resp, err := http.Get(url)
	if err != nil {
		logger.Error().Err(err).Str("keyword", keyword).Msg("Failed to search symbols from Alpha Vantage")
		return nil, NewAppError(http.StatusInternalServerError, types.StatusInternalServerError, "Failed to search symbols", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		logger.Error().Int("status", resp.StatusCode).Str("keyword", keyword).Msg("Alpha Vantage API returned non-200 status")
		return nil, NewAppError(http.StatusInternalServerError, types.StatusInternalServerError, "Alpha Vantage API error", nil)
	}

	var search SearchResponse
	if err := json.NewDecoder(resp.Body).Decode(&search); err != nil {
		logger.Error().Err(err).Str("keyword", keyword).Msg("Failed to parse search response")
		return nil, NewAppError(http.StatusInternalServerError, types.StatusInternalServerError, "Failed to parse search results", err)
	}

	if len(search.BestMatches) == 0 {
		logger.Warn().Str("keyword", keyword).Msg("No matching symbols found")
		return &search, nil // Empty results are valid
	}

	return &search, nil
}

// isValidInterval checks if the interval is supported by TIME_SERIES_INTRADAY
func isValidInterval(interval string) bool {
	validIntervals := []string{"1min", "5min", "15min", "30min", "60min"}
	for _, v := range validIntervals {
		if interval == v {
			return true
		}
	}
	return false
}

func SearchStockHandler(logger *zerolog.Logger) echo.HandlerFunc {
	return func(c echo.Context) error {
		keyword := c.QueryParam("query")
		if keyword == "" {
			logger.Error().Msg("Missing query parameter for stock search")
			return NewAppError(http.StatusBadRequest, types.StatusBadRequest, "Query parameter is required", nil)
		}

		searchResult, err := SearchSymbol(keyword, logger)
		if err != nil {
			return err // AppError already set
		}

		return c.JSON(http.StatusOK, map[string]interface{}{
			"message": types.StatusOK,
			"results": searchResult.BestMatches,
		})
	}
}

func GetStockQuoteHandler(logger *zerolog.Logger) echo.HandlerFunc {
	return func(c echo.Context) error {
		symbol := c.Param("symbol")
		if symbol == "" {
			logger.Error().Msg("Missing symbol parameter for stock quote")
			return NewAppError(http.StatusBadRequest, types.StatusBadRequest, "Symbol parameter is required", nil)
		}

		quote, err := FetchQuote(symbol, logger)
		if err != nil {
			return err // AppError already set
		}

		return c.JSON(http.StatusOK, map[string]interface{}{
			"message": types.StatusOK,
			"quote":   quote,
		})
	}
}

func GetIntradayDataHandler(logger *zerolog.Logger) echo.HandlerFunc {
	return func(c echo.Context) error {
		symbol := c.Param("symbol")
		interval := c.QueryParam("interval")
		if symbol == "" || interval == "" {
			logger.Error().Str("symbol", symbol).Str("interval", interval).Msg("Missing symbol or interval parameter for intraday data")
			return NewAppError(http.StatusBadRequest, types.StatusBadRequest, "Symbol and interval parameters are required", nil)
		}

		intraday, err := FetchIntraday(symbol, interval, logger)
		if err != nil {
			return err // AppError already set
		}

		return c.JSON(http.StatusOK, map[string]interface{}{
			"message": types.StatusOK,
			"data":    intraday,
		})
	}
}
