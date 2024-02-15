package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/fabulias/amaris-interview/internal/entities"
	"github.com/fabulias/amaris-interview/internal/store"
	"github.com/gin-gonic/gin"
)

const baseURL = "https://api.openweathermap.org/data/2.5/weather"

var _ IWeatherService = (*WeatherService)(nil)

// https://api.openweathermap.org/data/2.5/weather?q=London,uk&appid=4d1e87732e94ffd1804af4941ec7b355
type IWeatherService interface {
	GetWeatherByCity(ctx *gin.Context, city string) (*entities.WeatherResp, error)
	GetWeatherByCountry(ctx *gin.Context, countryCode string) (*entities.WeatherData, error)
	CreateCity(ctx *gin.Context, city *entities.City) (*entities.City, error)
	DeleteCity(ctx *gin.Context, city string) error
}

type WeatherService struct {
	store      store.WeatherStore
	httpClient *http.Client
	apiKey     string
}

func NewWeatherService(store store.WeatherStore, client *http.Client, apiKey string) *WeatherService {
	return &WeatherService{
		store:      store,
		httpClient: client,
		apiKey:     apiKey,
	}
}

func (w *WeatherService) fetchWeatherData(ctx *gin.Context, cityName string) (*entities.WeatherResp, error) {
	url := fmt.Sprintf("%s?q=%s&appid=%s&units=metric", baseURL, cityName, w.apiKey)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	resp, err := w.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	weatherResp := &entities.WeatherResp{}
	err = json.Unmarshal(responseBody, weatherResp)
	if err != nil {
		return nil, err
	}

	return weatherResp, nil
}

func (w *WeatherService) GetWeatherByCity(ctx *gin.Context, cityName string) (*entities.WeatherResp, error) {
	// TODO: add in handler ToLower cityName to avoid errors for searching, same for insert and other methods which receive city.
	city, err := w.store.GetCity(ctx, cityName)
	if err != nil {
		return nil, errors.New("city doesn't exist")
	}

	weatherResp, err := w.fetchWeatherData(ctx, cityName)
	if err != nil {
		return nil, err
	}

	wd := &entities.WeatherData{
		CityID: city.ID,
	}
	wd.CopyFromWeatherResp(weatherResp)

	_, err = w.store.UpdateWeatherData(ctx, wd)
	if err != nil {
		return nil, err
	}

	return weatherResp, nil
}
func (w *WeatherService) GetWeatherByCountry(ctx *gin.Context, countryCode string) (*entities.WeatherData, error) {
	results, err := w.store.GetWeatherDataByCountryCode(ctx, countryCode)
	if err != nil {
		return nil, err
	}

	lenResults := len(results)
	weatherDataResp := &entities.WeatherData{}
	for _, result := range results {
		weatherDataResp.AddFromWeatherResp(result)
	}

	weatherDataResp.DividerWeatherData(lenResults)

	return weatherDataResp, nil
}
func (w *WeatherService) CreateCity(ctx *gin.Context, city *entities.City) (*entities.City, error) {
	weatherResp, err := w.fetchWeatherData(ctx, city.Name)
	if err != nil {
		return nil, err
	}
	wd := &entities.WeatherData{}
	wd.CopyFromWeatherResp(weatherResp)

	insertedCity, _, err := w.store.CreateCityAndWeatherData(ctx, city, wd)
	if err != nil {
		return nil, err
	}

	return insertedCity, nil
}
func (w *WeatherService) DeleteCity(ctx *gin.Context, cityName string) error {
	city, err := w.store.GetCity(ctx, cityName)
	if err != nil {
		return err
	}
	if err = w.store.DeleteCityAndWeatherData(ctx, city.ID); err != nil {
		return err
	}
	return nil
}
