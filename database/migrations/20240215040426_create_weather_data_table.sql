-- +goose Up
-- +goose StatementBegin
CREATE TABLE weather_data (
  id SERIAL PRIMARY KEY,
  city_id INTEGER NOT NULL REFERENCES cities(id) ON DELETE CASCADE,
  temp NUMERIC,
  feels_like NUMERIC,
  temp_min NUMERIC,
  temp_max NUMERIC,
  pressure INTEGER,
  humidity INTEGER,
  sea_level INTEGER,
  grnd_level INTEGER,
  speed NUMERIC,
  deg INTEGER,
  gust NUMERIC
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE weather_data;
-- +goose StatementEnd
