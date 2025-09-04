package alphavantage

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/pratyush934/tradealpha/server/types"
	"github.com/pratyush934/tradealpha/server/util"

	"github.com/rs/zerolog"
)

const (
	AAlphaVantageBaseURL = "https://www.alphavantage.co/query"
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
	url := fmt.Sprintf("%s?function=GLOBAL_QUOTE&symbol=%s&apikey=%s", AAlphaVantageBaseURL, symbol, os.Getenv("ALPHA_VANTAGE_KEY"))
	resp, err := http.Get(url)
	if err != nil {
		logger.Error().Err(err).Str("symbol", symbol).Msg("Failed to fetch quote from Alpha Vantage")
		return nil, util.NewAppError(http.StatusInternalServerError, types.StatusInternalServerError, "Failed to fetch stock quote", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		logger.Error().Int("status", resp.StatusCode).Str("symbol", symbol).Msg("Alpha Vantage API returned non-200 status")
		return nil, util.NewAppError(http.StatusInternalServerError, types.StatusInternalServerError, "Alpha Vantage API error", nil)
	}

	var quote QuoteResponse
	if err := json.NewDecoder(resp.Body).Decode(&quote); err != nil {
		logger.Error().Err(err).Str("symbol", symbol).Msg("Failed to parse quote response")
		return nil, util.NewAppError(http.StatusInternalServerError, types.StatusInternalServerError, "Failed to parse stock quote", err)
	}

	if quote.GlobalQuote.Symbol == "" {
		logger.Error().Str("symbol", symbol).Msg("Invalid symbol or no data returned")
		return nil, util.NewAppError(http.StatusBadRequest, types.StatusBadRequest, "Invalid stock symbol", nil)
	}

	return &quote, nil
}

// FetchIntraday retrieves intraday time series data for a symbol
func FetchIntraday(symbol, interval string, logger *zerolog.Logger) (*IntradayResponse, error) {
	if !isValidInterval(interval) {
		logger.Error().Str("interval", interval).Msg("Invalid interval for intraday data")
		return nil, util.NewAppError(http.StatusBadRequest, types.StatusBadRequest, "Invalid interval", nil)
	}

	url := fmt.Sprintf("%s?function=TIME_SERIES_INTRADAY&symbol=%s&interval=%s&apikey=%s&extended_hours=true", AAlphaVantageBaseURL, symbol, interval, os.Getenv("ALPHA_VANTAGE_KEY"))
	resp, err := http.Get(url)
	if err != nil {
		logger.Error().Err(err).Str("symbol", symbol).Msg("Failed to fetch intraday data from Alpha Vantage")
		return nil, util.NewAppError(http.StatusInternalServerError, types.StatusInternalServerError, "Failed to fetch intraday data", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		logger.Error().Int("status", resp.StatusCode).Str("symbol", symbol).Msg("Alpha Vantage API returned non-200 status")
		return nil, util.NewAppError(http.StatusInternalServerError, types.StatusInternalServerError, "Alpha Vantage API error", nil)
	}

	var intraday IntradayResponse
	if err := json.NewDecoder(resp.Body).Decode(&intraday); err != nil {
		logger.Error().Err(err).Str("symbol", symbol).Msg("Failed to parse intraday response")
		return nil, util.NewAppError(http.StatusInternalServerError, types.StatusInternalServerError, "Failed to parse intraday data", err)
	}

	if intraday.MetaData.Symbol == "" {
		logger.Error().Str("symbol", symbol).Msg("Invalid symbol or no data returned")
		return nil, util.NewAppError(http.StatusBadRequest, types.StatusBadRequest, "Invalid stock symbol", nil)
	}

	return &intraday, nil
}

// SearchSymbol searches for stocks by keyword
func SearchSymbol(keyword string, logger *zerolog.Logger) (*SearchResponse, error) {
	// Check if API key is set
	apiKey := os.Getenv("ALPHA_VANTAGE_KEY")
	if apiKey == "" {
		logger.Error().Str("keyword", keyword).Msg("ALPHA_VANTAGE_KEY environment variable not set")
		return nil, util.NewAppError(http.StatusInternalServerError, types.StatusInternalServerError, "API key not configured", nil)
	}

	// Construct URL and log it (with API key masked)
	url := fmt.Sprintf("%s?function=SYMBOL_SEARCH&keywords=%s&apikey=%s", AAlphaVantageBaseURL, keyword, apiKey)
	logger.Info().Str("url", strings.Replace(url, apiKey, "****", -1)).Msg("Sending request to Alpha Vantage")

	// Make HTTP request
	resp, err := http.Get(url)
	if err != nil {
		logger.Error().Err(err).Str("keyword", keyword).Msg("Failed to search symbols")
		return nil, util.NewAppError(http.StatusInternalServerError, types.StatusInternalServerError, "Failed to search symbols", err)
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Error().Err(err).Str("keyword", keyword).Msg("Failed to read response body")
		return nil, util.NewAppError(http.StatusInternalServerError, types.StatusInternalServerError, "Failed to read API response", err)
	}
	logger.Info().Str("response", string(body)).Msg("Alpha Vantage raw response")

	// Check HTTP status code
	if resp.StatusCode != http.StatusOK {
		logger.Error().Int("status", resp.StatusCode).Str("keyword", keyword).Msg("Alpha Vantage API returned non-200 status")
		return nil, util.NewAppError(http.StatusInternalServerError, types.StatusInternalServerError, "Alpha Vantage API error", nil)
	}

	// Check for API error messages in JSON response
	var rawResponse map[string]interface{}
	if err := json.Unmarshal(body, &rawResponse); err != nil {
		logger.Error().Err(err).Str("keyword", keyword).Msg("Failed to parse raw response")
		return nil, util.NewAppError(http.StatusInternalServerError, types.StatusInternalServerError, "Failed to parse search results", err)
	}
	if info, exists := rawResponse["Information"]; exists {
		logger.Error().Str("keyword", keyword).Str("error", fmt.Sprintf("%v", info)).Msg("Alpha Vantage API error response")
		return nil, util.NewAppError(http.StatusBadGateway, types.StatusBadGateway, fmt.Sprintf("Alpha Vantage API error: %v", info), nil)
	}
	if note, exists := rawResponse["Note"]; exists {
		logger.Error().Str("keyword", keyword).Str("error", fmt.Sprintf("%v", note)).Msg("Alpha Vantage API rate limit exceeded")
		return nil, util.NewAppError(http.StatusTooManyRequests, types.StatusTooManyRequests, fmt.Sprintf("Alpha Vantage API rate limit: %v", note), nil)
	}

	// Parse response into SearchResponse struct
	var search SearchResponse
	if err := json.Unmarshal(body, &search); err != nil {
		logger.Error().Err(err).Str("keyword", keyword).Msg("Failed to parse search response")
		return nil, util.NewAppError(http.StatusInternalServerError, types.StatusInternalServerError, "Failed to parse search results", err)
	}

	// Log the number of matches found
	logger.Info().Int("matches", len(search.BestMatches)).Str("keyword", keyword).Msg("Parsed search results")

	if len(search.BestMatches) == 0 {
		logger.Warn().Str("keyword", keyword).Msg("No matching symbols found")
		return &search, nil
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
			return util.NewAppError(http.StatusBadRequest, types.StatusBadRequest, "Query parameter is required", nil)
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
			return util.NewAppError(http.StatusBadRequest, types.StatusBadRequest, "Symbol parameter is required", nil)
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
			return util.NewAppError(http.StatusBadRequest, types.StatusBadRequest, "Symbol and interval parameters are required", nil)
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

type DailyResponse struct {
	MetaData struct {
		Information   string `json:"1. Information"`
		Symbol        string `json:"2. Symbol"`
		LastRefreshed string `json:"3. Last Refreshed"`
		OutputSize    string `json:"4. Output Size"`
		TimeZone      string `json:"5. Time Zone"`
	} `json:"Meta Data"`
	TimeSeries map[string]struct {
		Open   string `json:"1. open"`
		High   string `json:"2. high"`
		Low    string `json:"3. low"`
		Close  string `json:"4. close"`
		Volume string `json:"5. volume"`
	} `json:"Time Series (Daily)"`
}

func GetDailyDataHandler(logger *zerolog.Logger) echo.HandlerFunc {
	return func(c echo.Context) error {
		symbol := c.Param("symbol")
		if symbol == "" {
			return util.NewAppError(http.StatusBadRequest, types.StatusBadRequest, "symbol is required", nil)
		}

		apiKey := os.Getenv("ALPHA_VANTAGE_KEY") // Replace with env variable or config
		url := fmt.Sprintf("https://www.alphavantage.co/query?function=TIME_SERIES_DAILY&symbol=%s&apikey=%s", symbol, apiKey)

		resp, err := http.Get(url)
		if err != nil {
			logger.Error().Err(err).Str("symbol", symbol).Msg("Failed to fetch daily data from Alpha Vantage")
			return util.NewAppError(http.StatusInternalServerError, types.StatusInternalServerError, "failed to fetch daily data", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			logger.Error().Str("symbol", symbol).Int("status", resp.StatusCode).Msg("Alpha Vantage API returned non-200 status")
			return util.NewAppError(http.StatusBadGateway, types.StatusBadGateway, "Alpha Vantage API error", nil)
		}

		var dailyData DailyResponse
		if err := json.NewDecoder(resp.Body).Decode(&dailyData); err != nil {
			logger.Error().Err(err).Str("symbol", symbol).Msg("Failed to parse daily data response")
			return util.NewAppError(http.StatusInternalServerError, types.StatusInternalServerError, "failed to parse daily data", err)
		}

		if _, ok := dailyData.TimeSeries["Error Message"]; ok {
			logger.Error().Str("symbol", symbol).Msg("Invalid symbol or API error")
			return util.NewAppError(http.StatusBadRequest, types.StatusBadRequest, "invalid symbol or API error", nil)
		}

		logger.Info().Str("symbol", symbol).Msg("Successfully fetched daily data")
		return c.JSON(http.StatusOK, dailyData)
	}
}
