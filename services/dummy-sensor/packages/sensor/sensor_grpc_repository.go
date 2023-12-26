package sensor

import "github.com/LukasJatmiko/simple-microservices/dummy-sensor/driver"

type SensorGRPCRepository struct {
	driver driver.Driver
}

func (sr *SensorGRPCRepository) SendDataCapture(data *SensorData) error {
	//to do
	sr.driver.GetInstance()
	return nil
}
