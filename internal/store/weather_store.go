package store

import (
	"context"
	"fmt"

	"github.com/Masterminds/squirrel"
	sq "github.com/Masterminds/squirrel"
	"github.com/fabulias/amaris-interview/internal/config"
	"github.com/fabulias/amaris-interview/internal/entities"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type WeatherStore interface {
	CreateCity(ctx *gin.Context, city *entities.City) (*entities.City, error)
	DeleteCityAndWeatherData(ctx *gin.Context, cityID int64) error
	GetCity(ctx *gin.Context, citySearch string) (*entities.City, error)

	CreateWeatherData(ctx *gin.Context, weatherData *entities.WeatherData) (*entities.WeatherData, error) // migrate this to trx way like Delete
	UpdateWeatherData(ctx *gin.Context, weatherData *entities.WeatherData) (*entities.WeatherData, error)
	GetWeatherDataByCountryCode(ctx context.Context, countryCode string) ([]*entities.WeatherData, error)
}

type WeatherPostgresStore struct {
	db *sqlx.DB
}

func NewWeatherPostgresStore(cfg config.Config) (WeatherStore, error) {
	dsn := fmt.Sprintf(
		"host=%s dbname=%s user=%s password=%s sslmode=disable",
		cfg.PostgresHost,
		cfg.PostgresDB,
		cfg.PostgresUser,
		cfg.PostgresPassword,
	)

	db, err := sqlx.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}

	return &WeatherPostgresStore{
		db: db,
	}, nil
}

func (w *WeatherPostgresStore) GetCity(ctx *gin.Context, citySearch string) (*entities.City, error) {
	query, args, err := sq.
		Select("*").
		From("cities").
		Where(sq.Eq{"name": citySearch}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, err
	}

	var city entities.City
	err = w.db.GetContext(ctx, &city, query, args...)
	if err != nil {
		return nil, err
	}
	return &city, nil
}

func (w *WeatherPostgresStore) CreateCity(ctx *gin.Context, city *entities.City) (*entities.City, error) {
	query, args, err := sq.
		Insert("cities").
		Columns("name", "country_code").
		Values(city.Name, city.CountryCode).
		Suffix("RETURNING *").
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, err
	}

	ans := entities.City{}
	err = w.db.QueryRowxContext(ctx, query, args...).StructScan(&ans)
	if err != nil {
		return nil, fmt.Errorf("error creating city: %v", err)
	}
	return &ans, nil
}

func (w *WeatherPostgresStore) CreateWeatherData(ctx *gin.Context, weatherData *entities.WeatherData) (*entities.WeatherData, error) {
	query, args, err := sq.
		Insert("weather_data").
		Columns(
			"city_id",
			"temp",
			"feels_like",
			"temp_min",
			"temp_max",
			"pressure",
			"humidity",
			"sea_level",
			"grnd_level",
			"speed",
			"deg",
			"gust",
		).
		Values(
			weatherData.CityID,
			weatherData.Temp,
			weatherData.FeelsLike,
			weatherData.TempMin,
			weatherData.TempMax,
			weatherData.Pressure,
			weatherData.Humidity,
			weatherData.SeaLevel,
			weatherData.GrndLevel,
			weatherData.WindSpeed,
			weatherData.WindDeg,
			weatherData.WindGust,
		).
		Suffix("RETURNING *").
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, err
	}

	ans := entities.WeatherData{}
	err = w.db.QueryRowxContext(ctx, query, args...).StructScan(&ans)
	if err != nil {
		return nil, fmt.Errorf("error creating weather data: %v", err)
	}
	return &ans, nil
}

func (w *WeatherPostgresStore) UpdateWeatherData(ctx *gin.Context, weatherData *entities.WeatherData) (*entities.WeatherData, error) {
	query, args, err := sq.
		Update("weather_data").
		SetMap(sq.Eq{
			"temp":       weatherData.Temp,
			"feels_like": weatherData.FeelsLike,
			"temp_min":   weatherData.TempMin,
			"temp_max":   weatherData.TempMax,
			"pressure":   weatherData.Pressure,
			"humidity":   weatherData.Humidity,
			"sea_level":  weatherData.SeaLevel,
			"grnd_level": weatherData.GrndLevel,
			"speed":      weatherData.WindSpeed,
			"deg":        weatherData.WindDeg,
			"gust":       weatherData.WindGust,
		}).
		Where(sq.Eq{"city_id": weatherData.CityID}).
		Suffix("RETURNING *").
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, err
	}

	ans := entities.WeatherData{}
	err = w.db.QueryRowxContext(ctx, query, args...).StructScan(&ans)
	if err != nil {
		return nil, fmt.Errorf("error updating weather data: %v", err)
	}
	return &ans, nil
}

func (w *WeatherPostgresStore) DeleteCityAndWeatherData(ctx *gin.Context, cityID int64) error {
	// Start a trx
	tx, err := w.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	_, err = squirrel.Delete("weather_data").
		Where("city_id = ?", cityID).
		PlaceholderFormat(sq.Dollar).
		RunWith(tx).Exec()
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = squirrel.Delete("cities").
		Where("id = ?", cityID).
		PlaceholderFormat(sq.Dollar).
		RunWith(tx).Exec()
	if err != nil {
		tx.Rollback()
		return err
	}

	// Commit the trx
	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (w *WeatherPostgresStore) GetWeatherDataByCountryCode(ctx context.Context, countryCode string) ([]*entities.WeatherData, error) {
	query, args, err := sq.
		Select("temp, feels_like, temp_min, temp_max, pressure, humidity, sea_level, grnd_level, speed, deg, gust").
		From("weather_data wd").
		Join("cities c ON wd.city_id = c.id").
		Where(sq.Eq{"c.country_code": countryCode}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := w.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var weatherData []*entities.WeatherData
	for rows.Next() {
		var wd entities.WeatherData
		// Scan into the WeatherData struct fields
		if err := rows.Scan(&wd.Temp, &wd.FeelsLike, &wd.TempMin, &wd.TempMax, &wd.Pressure, &wd.Humidity, &wd.SeaLevel, &wd.GrndLevel, &wd.WindSpeed, &wd.WindDeg, &wd.WindGust); err != nil {
			return nil, err
		}
		weatherData = append(weatherData, &wd)
	}

	return weatherData, nil
}
