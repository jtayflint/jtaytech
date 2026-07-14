package integrations

import (
	"context"
	"regexp"
)

// Regex to strictly enforce max 4 decimal places for both negative and positive coordinates
var CoordRegex = regexp.MustCompile(`^-?\d+(\.\d{1,4})?$`)

type WeatherResponse struct {
	ShortForecast string `json:"short_forecast"`
	Temperature   int    `json:"temperature"`
	Unit          string `json:"unit"`
	Category      string `json:"category"` // hot, cold, or moderate
}

// WeatherService defines the interface for interacting with the weather service provider.
type WeatherService interface {
	GetForecast(ctx context.Context, lat, lon, userAgent string) (*WeatherResponse, error)
}
