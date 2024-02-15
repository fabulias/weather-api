package entities

type City struct {
	ID          int64  `db:"id,primarykey,autoIncrement"`
	Name        string `json:"name" db:"name"`
	CountryCode string `json:"countryCode" db:"country_code"`
}
