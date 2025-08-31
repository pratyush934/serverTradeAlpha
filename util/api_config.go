package util

import "os"

const (
	AlphaBaseURl = "https://www.alphavantage.co/query"
)

var APIKey = os.Getenv("ALPHA_VANTAGE_KEY")
