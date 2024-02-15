package app

import (
	"fmt"
	"net/http"
	"time"

	"github.com/fabulias/amaris-interview/internal/config"
	"github.com/fabulias/amaris-interview/internal/handlers"
	"github.com/fabulias/amaris-interview/internal/middleware"
	"github.com/fabulias/amaris-interview/internal/service"
	"github.com/fabulias/amaris-interview/internal/store"
	"go.uber.org/zap"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type App struct {
	server          *gin.Engine
	port            string
	weatherHandlers *handlers.WeatherHandlers
}

func NewServer(cfg config.Config, logger *zap.Logger) (*App, error) {
	r := gin.Default()
	r.Use(cors.Default())
	r.Use(middleware.LoggingMiddleware(logger))

	v1Group := r.Group("v1")

	postgresStore, err := store.NewWeatherPostgresStore(cfg)
	if err != nil {
		return nil, err
	}
	httpClient := &http.Client{
		Timeout: 1 * time.Minute,
	}

	weatherService := service.NewWeatherService(postgresStore, httpClient, cfg.ApiKey)
	weatherHandlers := handlers.NewWeatherHandlers(weatherService)

	v1Group.GET("/health", weatherHandlers.GetHealth)
	v1Group.GET("/weather/city/:city", weatherHandlers.GetWeatherByCity)
	v1Group.GET("/weather/country/:country", weatherHandlers.GetWeatherByCountry)
	v1Group.POST("/weather/city", weatherHandlers.CreateCity)
	v1Group.DELETE("/weather/city/:city", weatherHandlers.DeleteCity)

	return &App{
		server:          r,
		weatherHandlers: weatherHandlers,
		port:            fmt.Sprintf(":%s", cfg.ServicePort),
	}, nil
}

func (app *App) Run() error {
	return app.server.Run(app.port)
}
