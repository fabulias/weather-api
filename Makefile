run:
	@echo "running app..."
	@docker-compose up --build

down:
	@echo "shuting down app..."
	@docker-compose down

run-migrations:
	@echo "running migrations..."
	@goose -dir database/migrations postgres "host=localhost user=postgres password=password dbname=weather_api sslmode=disable" up
