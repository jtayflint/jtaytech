package web

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"weatherapi/app"
	"weatherapi/integrations"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type WeatherHandler struct {
	service integrations.WeatherService
	logger  *zap.Logger
	app     *app.App
}

func NewWeatherHandler(app *app.App, logger *zap.Logger, svc integrations.WeatherService) *WeatherHandler {
	return &WeatherHandler{
		app:     app,
		logger:  logger,
		service: svc,
	}
}

func (h *WeatherHandler) GetWeather(w http.ResponseWriter, r *http.Request) {
	lat := chi.URLParam(r, "lat")
	lon := chi.URLParam(r, "lon")

	// Validate coordinates exist and conform to 4 decimal places max
	if !integrations.CoordRegex.MatchString(lat) || !integrations.CoordRegex.MatchString(lon) {
		h.logger.Warn("Failed request: Coordinate validation failed", zap.String("lat", lat), zap.String("lon", lon))
		http.Error(w, "Invalid coordinate format. Coordinates limited to 4 decimal places.", http.StatusBadRequest)
		return
	}

	// Range validation checks
	latVal, _ := strconv.ParseFloat(lat, 64)
	lonVal, _ := strconv.ParseFloat(lon, 64)
	if latVal < -90 || latVal > 90 || lonVal < -180 || lonVal > 180 {
		h.logger.Warn("Failed request: Coordinates out of bounds", zap.Float64("lat", latVal), zap.Float64("lon", lonVal))
		http.Error(w, "Invalid coordinates. Latitude must be between -90/90 and Longitude between -180/180.", http.StatusBadRequest)
		return
	}

	// Establish correct User-Agent string
	ua := r.Header.Get("User-Agent")
	if ua == "" {
		ua = h.app.DefaultUserAgent
	}

	// Interface call
	weather, err := h.service.GetForecast(r.Context(), lat, lon, ua)
	if err != nil {
		h.logger.Error("Failed request: NWS integration error", zap.Error(err), zap.String("lat", lat), zap.String("lon", lon))
		http.Error(w, fmt.Sprintf("Weather service error: %v", err), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode(weather); err != nil {
		h.logger.Error("failed to encode weather response: %v", zap.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func NewRouter(lc fx.Lifecycle, shutdowner fx.Shutdowner, app *app.App, wh *WeatherHandler) {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	// RESTful path parameters layout
	r.Get("/weather/lat/{lat}/lon/{lon}", wh.GetWeather)
	srv := &http.Server{
		Addr:    ":" + app.Port,
		Handler: r,
	}
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {

			wh.logger.Info("Starting HTTP Server", zap.String("port", app.Port))
			go func() {
				if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
					wh.logger.Error("HTTP server failed", zap.Error(err))
					if inerr := shutdowner.Shutdown(); inerr != nil {
						wh.logger.Error("Failed to shutdown application", zap.Error(inerr))
					}
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			wh.logger.Info("Stopping HTTP Server")
			return srv.Shutdown(ctx)
		},
	})
}
