package sensor

import (
	"fmt"
	"time"

	"github.com/LukasJatmiko/simple-microservices/dummy-sensor/driver"
)

type SensorData struct {
	Value     float64   `json:"value"`
	Type      string    `json:"type"`
	ID1       string    `json:"id1"`
	ID2       int       `json:"id2"`
	Timestamp time.Time `json:"timestamp"`
}

type SensorHandler interface {
	SendDataCapture(data *SensorData) error
}

func NewSensorHandler(drv driver.Driver) (SensorHandler, error) {
	switch drv.GetOptions().Type {
	case driver.GRPC:
		//to do
		return nil, fmt.Errorf("unknown driver type")
	case driver.RABBITMQ:
		return &SensorRMQRepository{
			driver: drv,
		}, nil
	default:
		return nil, fmt.Errorf("unknown driver type")
	}
}
