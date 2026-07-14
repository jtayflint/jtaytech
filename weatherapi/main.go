package main

import (
	"weatherapi/app"
	"weatherapi/nws"
	"weatherapi/web"

	"go.uber.org/fx"
	"go.uber.org/zap"
)

// ==========================================
// 6. Main Entrypoint
// ==========================================

func main() {
	fx.New(
		fx.Provide(
			zap.NewDevelopment, // Structed logging via Zap
			app.NewApp,
			nws.NewWeatherService,
			web.NewWeatherHandler,
		),
		fx.Invoke(web.NewRouter),
	).Run()
}
