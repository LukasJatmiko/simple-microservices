package main

import (
	"fmt"
	"time"

	"github.com/LukasJatmiko/simple-microservices/data-processor/constants"
	"github.com/LukasJatmiko/simple-microservices/data-processor/driver"
	"github.com/LukasJatmiko/simple-microservices/data-processor/packages/api"
	"github.com/LukasJatmiko/simple-microservices/data-processor/packages/sensor"
	"github.com/LukasJatmiko/simple-microservices/data-processor/utils"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// Echo instance
	ec := echo.New()

	rabbitmqDriver := driver.NewDriver(&driver.Options{
		Type: driver.RABBITMQ,
		Rabbitmq: driver.RabbitmqOptions{
			URL:          utils.GetEnvOrDefaultString(string(constants.ENVRabbitmqURL), "amqp://guest:guest@localhost:5672/"),
			ExchangeName: utils.GetEnvOrDefaultString(string(constants.ENVRabbitmqXchName), "simple-ms"),
			ExchangeType: driver.Topic,
			QueueName:    utils.GetEnvOrDefaultString(string(constants.ENVRabbitmqQName), "q-sensor-data-stream"),
			RoutingKey:   utils.GetEnvOrDefaultString(string(constants.ENVRabbitmqRoutingKey), "sensor-data-stream"),
			ConsumerName: utils.GetEnvOrDefaultString(string(constants.ENVRabbitmqConsumerName), "data-processor-svc"),
		},
	})
	if e := rabbitmqDriver.Init(); e != nil {
		ec.Logger.Fatal(e)
	}

	var gormDialector gorm.Dialector
	switch driver.Database(utils.GetEnvOrDefaultString(string(constants.ENVDBDriverType), string(driver.MYSQL))) {
	case driver.POSTGRES:
		gormDialector = postgres.Open(utils.GetEnvOrDefaultString(string(constants.ENVDBURI), "postgresql://postgres:postgres@localhost:5432/postgres?sslmode=disabled"))
	default:
		gormDialector = mysql.Open(utils.GetEnvOrDefaultString(string(constants.ENVDBURI), "mysql://user:password@localhost:3306/mysql"))
	}
	gormDriver := driver.NewDriver(&driver.Options{
		Type: driver.GORM,
		Gorm: driver.GormOptions{
			Dialector:             gormDialector,
			MaxOpenConnection:     utils.GetEnvOrDefaultInt(string(constants.ENVMaxOpenConn), 5),
			MaxIdleConnection:     utils.GetEnvOrDefaultInt(string(constants.ENVMaxIdleConn), 1),
			MaxConnectionLifetime: utils.GetEnvOrDefaultDuration(string(constants.ENVMaxConnLifetime), (1800 * time.Second)),
		},
	})

	if e := gormDriver.Init(); e != nil {
		ec.Logger.Fatal(e)
	}

	if sensorHandler, e := sensor.NewSensorhandler(rabbitmqDriver, gormDriver, ec); e != nil {
		ec.Logger.Fatal(e)
	} else {
		go sensorHandler.HandleSensorData()
	}

	// Middleware
	ec.Use(middleware.Logger())
	ec.Use(middleware.Recover())

	// Routes
	RGapi := ec.Group("/api")
	RGv1 := RGapi.Group("/v1")
	RGsensor := RGv1.Group("/sensor")

	apiHandler := api.NewAPIHandler(gormDriver)
	RGsensor.GET("/", apiHandler.GetSensorData(new(api.APIFilter)))
	RGsensor.DELETE("/", apiHandler.DeleteSensorData(new(api.APIFilter)))
	RGsensor.PATCH("/", apiHandler.UpdateSensorData(new(api.APIFilter)))

	// Start server
	appPort := utils.GetEnvOrDefaultString(string(constants.ENVAppPort), "8081")
	ec.Logger.Fatal(ec.Start(fmt.Sprintf(":%v", appPort)))
}
