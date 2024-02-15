# weather-api

This project is written in Golang using Docker-compose and Goose

## How to run Weather API implementation

You need to have installed [Docker-compose](https://docs.docker.com/compose/), [Docker](https://www.docker.com/) and [Goose](https://github.com/pressly/goose) in your machine.

1. Download the project and be sure you are in the root of the project.
2. Run `make run` command.
3. Run `make run-migrations` command.
4. Enjoy!

## Weather API v1 methods 

**Base URL:** `http://localhost:9000/v1`

**Endpoints:**

| Method | Path | Description | Request Body | Response |
|---|---|---|---|---|
| `GET` | `/health` | Check API health | - | `200 OK: {}` |
| `GET` | `/weather/city/:city` | Get weather for a city | - | `200 OK: { // data information from OpenWeatherMap API }` |
| `GET` | `/weather/country/:country` | Get weather for a country | - | `200 OK: {// average data information from OpenWeatherMap API ` |
| `POST` | `/weather/city` | Create a new city | `{ "name": "New York" }` | `201 Created: {"ID":1,"name":"santiago","countryCode":"uk"}` |
| `DELETE` | `/weather/city/:city` | Delete a city | - | `204 No Content` |

**Notes:**

- Error responses follow standard HTTP status codes with JSON error messages.
- Authentication and authorization might be required for certain endpoints.


## Testing

You can download the collection from `Postman` folder (Collection and Environment variables).

## Pending topics to cover

* Create unit and integration tests.
* Return better errors, at this moment it's only returning errors created by standard library of Golang and from used packages.
* It was not mandatory, but with more time we should remove all the environment variables from this repository and inject directly as secrets to containers.