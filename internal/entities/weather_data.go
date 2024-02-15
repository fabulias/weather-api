package entities

type WeatherData struct {
	ID        int64   `db:"id"`
	CityID    int64   `db:"city_id"`
	Temp      float64 `db:"temp"`
	FeelsLike float64 `db:"feels_like"`
	TempMin   float64 `db:"temp_min"`
	TempMax   float64 `db:"temp_max"`
	Pressure  int     `db:"pressure"`
	Humidity  int     `db:"humidity"`
	SeaLevel  int     `db:"sea_level"`
	GrndLevel int     `db:"grnd_level"`
	WindSpeed float64 `db:"speed"`
	WindDeg   int     `db:"deg"`
	WindGust  float64 `db:"gust"`
}

// Function to copy values from WeatherResp to WeatherData
func (wd *WeatherData) CopyFromWeatherResp(resp *WeatherResp) {
	// Copy relevant values from Main struct
	wd.Temp = resp.Main.Temp
	wd.FeelsLike = resp.Main.FeelsLike
	wd.TempMin = resp.Main.TempMin
	wd.TempMax = resp.Main.TempMax
	wd.Pressure = resp.Main.Pressure
	wd.Humidity = resp.Main.Humidity

	// Copy relevant values from Wind struct
	wd.WindSpeed = resp.Wind.Speed
	wd.WindDeg = resp.Wind.Deg
	wd.WindGust = resp.Wind.Gust
}

// Function to add values from WeatherData input
func (wd *WeatherData) AddFromWeatherResp(input *WeatherData) {
	wd.Temp += input.Temp
	wd.FeelsLike += input.FeelsLike
	wd.TempMin += input.TempMin
	wd.TempMax += input.TempMax
	wd.Pressure += input.Pressure
	wd.Humidity += input.Humidity

	wd.WindSpeed += input.WindSpeed
	wd.WindDeg += input.WindDeg
	wd.WindGust += input.WindGust
}

// Function to add values from WeatherData input
func (wd *WeatherData) DividerWeatherData(divider int) {
	wd.Temp /= float64(divider)
	wd.FeelsLike /= float64(divider)
	wd.TempMin /= float64(divider)
	wd.TempMax /= float64(divider)
	wd.Pressure /= divider
	wd.Humidity /= divider

	wd.WindSpeed /= float64(divider)
	wd.WindDeg /= divider
	wd.WindGust /= float64(divider)
}
