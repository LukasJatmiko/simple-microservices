package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/LukasJatmiko/simple-microservices/dummy-sensor/constants"
	"github.com/LukasJatmiko/simple-microservices/dummy-sensor/driver"
	"github.com/LukasJatmiko/simple-microservices/dummy-sensor/packages/api"
	"github.com/LukasJatmiko/simple-microservices/dummy-sensor/packages/sensor"
	"github.com/LukasJatmiko/simple-microservices/dummy-sensor/utils"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	rabbitmqDriver := driver.NewDriver(&driver.Options{
		Type: driver.RABBITMQ,
		Rabbitmq: driver.RabbitmqOptions{
			URL:          utils.GetEnvOrDefaultString(string(constants.ENVRabbitmqURL), "amqp://guest:guest@localhost:5672/"),
			ExchangeName: utils.GetEnvOrDefaultString(string(constants.ENVRabbitmqXchName), "simple-ms"),
			ExchangeType: driver.Topic,
			QueueName:    utils.GetEnvOrDefaultString(string(constants.ENVRabbitmqQName), "q-sensor-data-stream"),
			RoutingKey:   utils.GetEnvOrDefaultString(string(constants.ENVRabbitmqRoutingKey), "sensor-data-stream"),
		},
	})
	if e := rabbitmqDriver.Init(); e != nil {
		panic(e)
	}

	durationUpdater := make(chan time.Duration)
	duration := (10 * time.Second)

	if sensorHandler, e := sensor.NewSensorHandler(rabbitmqDriver); e != nil {
		panic(e)
	} else {
		go func() {
			for duration = range durationUpdater {
				//
			}
		}()
		go func() {
			for {
				data := sensor.SensorData{
					Type:      utils.GetEnvOrDefaultString(string(constants.ENVDataSensorType), "TEMPERATURE"),
					ID1:       utils.GetEnvOrDefaultString(string(constants.ENVDataSensorID1), "FO"),
					ID2:       utils.GetEnvOrDefaultInt(string(constants.ENVDataSensorID2), 1),
					Value:     float64(randomNumber(20, 28)),
					Timestamp: time.Now(),
				}
				sensorHandler.SendDataCapture(&data)
				time.Sleep(duration)
			}
		}()
	}

	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Routes
	apiHandler := api.NewApiHandler(&api.Options{ChSensorDataInterval: durationUpdater})
	e.POST("/update", apiHandler.UpdateSensorDataInterval())

	// Start server
	appPort := utils.GetEnvOrDefaultString(string(constants.ENVAppPort), "8080")
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%v", appPort)))
}

func randomNumber(min, max int) int {
	return (rand.Intn(max-min) + min)
}
