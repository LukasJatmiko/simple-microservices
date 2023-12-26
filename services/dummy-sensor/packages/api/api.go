package api

import (
	"time"

	"github.com/labstack/echo/v4"
)

type APIRequestPayload struct {
	SensorConfiguration SensorConfiguration `json:"sensor_configuration"`
}

type APIResponsePayload struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

type SensorConfiguration struct {
	SensorDatainterval string `json:"sensor_data_interval"`
}

type APIHandler interface {
	UpdateSensorDataInterval() func(c echo.Context) error
}

type Options struct {
	ChSensorDataInterval chan time.Duration
}

func NewApiHandler(opt *Options) APIHandler {
	return &APIRepository{
		Options: opt,
	}
}
