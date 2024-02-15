package handlers

import (
	"net/http"

	"github.com/fabulias/amaris-interview/internal/entities"
	"github.com/fabulias/amaris-interview/internal/service"
	"github.com/gin-gonic/gin"
)

type WeatherHandlers struct {
	service service.IWeatherService
}

func NewWeatherHandlers(weatherService service.IWeatherService) *WeatherHandlers {
	return &WeatherHandlers{
		service: weatherService,
	}
}

func (w *WeatherHandlers) GetHealth(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{})
}

func (w *WeatherHandlers) GetWeatherByCity(c *gin.Context) {
	city := c.Param("city")
	resp, err := w.service.GetWeatherByCity(c, city)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)

}
func (w *WeatherHandlers) GetWeatherByCountry(c *gin.Context) {
	countryCode := c.Param("country")
	resp, err := w.service.GetWeatherByCountry(c, countryCode)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"countryCode": countryCode, "data": resp})
}

func (w *WeatherHandlers) CreateCity(c *gin.Context) {
	city := &entities.City{}
	if err := c.BindJSON(&city); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	cityResp, err := w.service.CreateCity(c, city)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	c.JSON(http.StatusCreated, cityResp)
}
func (w *WeatherHandlers) DeleteCity(c *gin.Context) {
	city := c.Param("city")
	if err := w.service.DeleteCity(c, city); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusNoContent, gin.H{})
}
