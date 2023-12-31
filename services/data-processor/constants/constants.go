package constants

import (
	"github.com/LukasJatmiko/simple-microservices/data-processor/types"
)

const AuthTypeStatic types.AuthType = "STATIC"
const AuthTypeJWT types.AuthType = "JWT"

/*
environments
*/
const ENVAppPort types.Environment = "APP_PORT"

const ENVDBURI types.Environment = "DB_URI"
const ENVDBDriverType types.Environment = "DB_DRIVER_TYPE"
const ENVMaxOpenConn types.Environment = "DB_MAX_OPEN_CONN"
const ENVMaxConnLifetime types.Environment = "DB_MAX_CONN_LIFETIME"
const ENVMaxIdleConn types.Environment = "DB_MAX_IDLE_CONN"

const ENVRabbitmqURL types.Environment = "RABBITMQ_URL"
const ENVRabbitmqXchName types.Environment = "RABBITMQ_XCH_NAME"
const ENVRabbitmqQName types.Environment = "RABBITMQ_Q_NAME"
const ENVRabbitmqRoutingKey types.Environment = "RABBITMQ_ROUTINGKEY"
const ENVRabbitmqConsumerName types.Environment = "RABBITMQ_CONSUMER_NAME"

const ENVAuthType types.Environment = "AUTH_TYPE"
const ENVAuthTokens types.Environment = "AUTH_TOKENS"
const ENVAuthJWTPrivateKey types.Environment = "AUTH_JWT_PRIVATE_KEY"
const ENVAuthJWTPublicKey types.Environment = "AUTH_JWT_PUBLIC_KEY"
const ENVAuthJWTExpiration types.Environment = "AUTH_JWT_EXPIRATION"
