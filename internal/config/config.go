package config

type Config struct {
	Environment          string `env:"ENVIRONMENT,required=true"`
	ServicePort          string `env:"SERVICE_PORT,required=true"`
	ServiceDatabaseTable string `env:"SERVICE_DATABASE_TABLE,required=true"`
	ApiKey               string `env:"API_KEY,required=true"`
	Database
}

type Database struct {
	PostgresHost     string `env:"POSTGRES_HOST,required=true"`
	PostgresDB       string `env:"POSTGRES_DB,required=true"`
	PostgresUser     string `env:"POSTGRES_USER,required=true"`
	PostgresPassword string `env:"POSTGRES_PASSWORD,required=true"`
}
