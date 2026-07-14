package app

import (
	"weatherapi/tools"

	"github.com/caarlos0/env/v11"
)

type Config struct {
	Port             string `env:"SERVER_PORT" envDefault:"8080"`
	DefaultUserAgent string `env:"DEFAULT_USER_AGENT" envDefault:"MyWeatherApp/1.0 (deepcreektech@gmail.com)"`
	ColdTemp         string `env:"COLD_TEMP" envDefault:"50"`
	HotTemp          string `env:"HOT_TEMP" envDefault:"85"`
}

type App struct {
	Port             string
	DefaultUserAgent string
	ColdTemp         float64
	HotTemp          float64
}

func NewApp() (*App, error) {
	var cfg Config
	if err := env.Parse(&cfg); err != nil {
		return nil, err
	}
	coldTemp, err := tools.NumFromString(cfg.ColdTemp)
	if err != nil {
		return nil, err
	}
	hotTemp, err := tools.NumFromString(cfg.HotTemp)
	if err != nil {
		return nil, err
	}
	return &App{
		Port:             cfg.Port,
		DefaultUserAgent: cfg.DefaultUserAgent,
		ColdTemp:         coldTemp,
		HotTemp:          hotTemp,
	}, nil
}
