package sensor

import (
	"context"
	"encoding/json"
	"time"

	"github.com/LukasJatmiko/simple-microservices/dummy-sensor/driver"
	amqp "github.com/rabbitmq/amqp091-go"
)

type SensorRMQRepository struct {
	driver driver.Driver
}

func (sr *SensorRMQRepository) SendDataCapture(data *SensorData) error {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if body, e := json.Marshal(data); e != nil {
		return e
	} else {
		sr.driver.GetInstance().(*driver.RabbitmqInstance).Channel.PublishWithContext(ctx,
			sr.driver.GetOptions().Rabbitmq.ExchangeName, // exchange
			sr.driver.GetOptions().Rabbitmq.RoutingKey,   // routing key
			false, // mandatory
			false, // immediate
			amqp.Publishing{
				ContentType: "application/json",
				Body:        body,
			},
		)
		return nil
	}
}
