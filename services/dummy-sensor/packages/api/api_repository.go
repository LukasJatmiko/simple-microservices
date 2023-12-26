package api

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

type APIRepository struct {
	Options *Options
}

func (ar *APIRepository) UpdateSensorDataInterval() func(c echo.Context) error {
	return (func(c echo.Context) error {
		payload := new(APIRequestPayload)
		if e := c.Bind(payload); e != nil {
			return c.JSON(http.StatusOK, APIResponsePayload{
				Status:  http.StatusBadRequest,
				Message: "payload format error",
				Data:    echo.Map{},
			})
		} else {
			if d, e := time.ParseDuration(payload.SensorConfiguration.SensorDatainterval); e != nil {
				return c.JSON(http.StatusOK, APIResponsePayload{
					Status:  http.StatusBadRequest,
					Message: "payload format error",
					Data:    echo.Map{},
				})
			} else {
				ar.Options.ChSensorDataInterval <- d
				return c.JSON(http.StatusOK, APIResponsePayload{
					Status:  http.StatusOK,
					Message: "sensor data sending interval updated",
					Data:    echo.Map{},
				})
			}
		}
	})
}
