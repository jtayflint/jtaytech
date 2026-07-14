package nws

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/caarlos0/env/v11"

	"weatherapi/app"
	"weatherapi/integrations"

	"go.uber.org/zap"
)

type nwsConfig struct {
	BaseURL    *string `env:"NWS_BASE_URL" envDefault:"https://api.weather.gov"`
	PointsPath *string `env:"NWS_POINTS_PATH" envDefault:"/points/%s,%s"`
}

// NWS metadata response structure
type nwsPointsResponse struct {
	Properties struct {
		Forecast string `json:"forecast"`
	} `json:"properties"`
}

// NWS forecast grid response structure
type nwsForecastResponse struct {
	Properties struct {
		Periods []struct {
			Number          int    `json:"number"`
			Name            string `json:"name"`
			Temperature     int    `json:"temperature"`
			TemperatureUnit string `json:"temperatureUnit"`
			ShortForecast   string `json:"shortForecast"`
			IsDaytime       bool   `json:"isDaytime"`
		} `json:"periods"`
	} `json:"properties"`
}

var nwsCfg nwsConfig

type nwsService struct {
	app        *app.App
	httpClient *http.Client
	logger     *zap.Logger
}

func NewWeatherService(app *app.App, logger *zap.Logger) (integrations.WeatherService, error) {

	if err := env.Parse(&nwsCfg); err != nil {
		return nil, err
	}
	return &nwsService{
		app: app,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
		logger: logger,
	}, nil
}

// Helper to characterize temperature ranges
func (s *nwsService) characterizeTemp(temp int, unit string) string {
	// Normalize to Fahrenheit for categorization standard
	fTemp := float64(temp)
	if strings.ToUpper(unit) == "C" {
		fTemp = (fTemp * 9 / 5) + 32
	}

	switch {
	case fTemp < s.app.ColdTemp:
		return "cold"
	case fTemp > s.app.HotTemp:
		return "hot"
	default:
		return "moderate"
	}
}

func (s *nwsService) GetForecast(ctx context.Context, lat, lon string, userAgent string) (*integrations.WeatherResponse, error) {
	// 1. Fetch metadata to get the grid forecast URL
	pointsURL := fmt.Sprintf("%s/points/%s,%s", *nwsCfg.BaseURL, lat, lon)

	req, err := http.NewRequestWithContext(ctx, "GET", pointsURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", userAgent)

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to call NWS points endpoint: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("NWS points endpoint returned status %d", resp.StatusCode)
	}

	var pointsData nwsPointsResponse
	if err := json.NewDecoder(resp.Body).Decode(&pointsData); err != nil {
		return nil, fmt.Errorf("failed to decode points response: %w", err)
	}

	// 2. Fetch the actual grid forecast using the URL returned from step 1
	forecastURL := pointsData.Properties.Forecast
	if forecastURL == "" {
		return nil, fmt.Errorf("grid forecast URL empty in NWS metadata")
	}

	req, err = http.NewRequestWithContext(ctx, "GET", forecastURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", userAgent)

	resp, err = s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to call NWS forecast endpoint: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("NWS forecast endpoint returned status %d", resp.StatusCode)
	}

	var forecastData nwsForecastResponse
	if err := json.NewDecoder(resp.Body).Decode(&forecastData); err != nil {
		return nil, fmt.Errorf("failed to decode forecast response: %w", err)
	}

	if len(forecastData.Properties.Periods) == 0 {
		return nil, fmt.Errorf("no forecast periods returned from NWS")
	}

	// Grab the first period (current/today)
	currentPeriod := forecastData.Properties.Periods[0]

	return &integrations.WeatherResponse{
		ShortForecast: currentPeriod.ShortForecast,
		Temperature:   currentPeriod.Temperature,
		Unit:          currentPeriod.TemperatureUnit,
		Category:      s.characterizeTemp(currentPeriod.Temperature, currentPeriod.TemperatureUnit),
	}, nil
}
