package sensor

import (
	"fmt"
	"time"

	"github.com/LukasJatmiko/simple-microservices/data-processor/driver"
	"github.com/labstack/echo/v4"
)

type SensorData struct {
	Value     float64   `json:"value" gorm:"column:value"`
	Type      string    `json:"type" gorm:"column:type"`
	ID1       string    `json:"id1" gorm:"column:id1"`
	ID2       int       `json:"id2" gorm:"column:id2"`
	Timestamp time.Time `json:"timestamp" gorm:"column:timestamp"`
}

type SensorHandler interface {
	HandleSensorData() error
}

func NewSensorhandler(messageDrv, databaseDrv driver.Driver, echo *echo.Echo) (SensorHandler, error) {
	switch messageDrv.GetOptions().Type {
	case driver.GRPC:
		//to do
		return nil, fmt.Errorf("unknown driver type")
	case driver.RABBITMQ:
		return &RabbitmqSensorHandler{
			echo:           echo,
			MessageDriver:  messageDrv,
			DatabaseDriver: databaseDrv,
		}, nil
	default:
		return nil, fmt.Errorf("unknown driver type")
	}
}
