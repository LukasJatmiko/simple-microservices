package api

import (
	"time"

	"github.com/LukasJatmiko/simple-microservices/data-processor/driver"
	"github.com/labstack/echo/v4"
)

type APIRequestPayload struct {
}

type APIResponsePayload struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

type APIFilter struct {
	ID1        string
	ID2        string
	start      time.Time
	end        time.Time
	Page       int
	RowPerPage int
	Order      string
}

type APIHandler interface {
	GetSensorData(f *APIFilter) func(c echo.Context) error
	DeleteSensorData(f *APIFilter) func(c echo.Context) error
	UpdateSensorData(f *APIFilter) func(c echo.Context) error
}

func NewAPIHandler(dbDriver driver.Driver) APIHandler {
	return &ApiRepository{
		DBDriver: dbDriver,
	}
}
