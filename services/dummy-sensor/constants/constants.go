package constants

import (
	"github.com/LukasJatmiko/simple-microservices/dummy-sensor/types"
)

/*
environments
*/
const ENVAppPort types.Environment = "APP_PORT"

const ENVRabbitmqURL types.Environment = "RABBITMQ_URL"
const ENVRabbitmqXchName types.Environment = "RABBITMQ_XCH_NAME"
const ENVRabbitmqQName types.Environment = "RABBITMQ_Q_NAME"
const ENVRabbitmqRoutingKey types.Environment = "RABBITMQ_ROUTINGKEY"

const ENVDataSensorType types.Environment = "DATA_SENSOR_TYPE"
const ENVDataSensorID1 types.Environment = "DATA_SENSOR_ID1"
const ENVDataSensorID2 types.Environment = "DATA_SENSOR_ID2"
