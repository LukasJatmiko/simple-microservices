package sensor

import (
	"encoding/json"

	"github.com/LukasJatmiko/simple-microservices/data-processor/driver"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type RabbitmqSensorHandler struct {
	echo           *echo.Echo
	MessageDriver  driver.Driver
	DatabaseDriver driver.Driver
}

func (rsh *RabbitmqSensorHandler) HandleSensorData() error {
	if messages, e := rsh.MessageDriver.GetInstance().(*driver.RabbitmqInstance).Channel.Consume(
		rsh.MessageDriver.GetOptions().Rabbitmq.QueueName,
		rsh.MessageDriver.GetOptions().Rabbitmq.ConsumerName,
		false,
		false,
		false,
		false,
		nil,
	); e != nil {
		return e
	} else {
		for message := range messages {
			sensorData := new(SensorData)
			if e := json.Unmarshal(message.Body, sensorData); e != nil {
				rsh.echo.Logger.Error(e)
			} else {
				//store sensor data to database
				db := rsh.DatabaseDriver.GetWrapperInstance().(*gorm.DB)
				result := db.Create(sensorData)
				if result.Error != nil {
					rsh.echo.Logger.Error(e)
				} else {
					//acknowledge message
					message.Ack(false)
				}
			}
		}
		return nil
	}
}
