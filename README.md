# simple-microservices

## System Diagram

![System Diagram](resources/system-diagram.png "System Diagram")



## Table structure

![System Diagram](resources/tables.png "System Diagram")


## Application directory structure
    .
    ├── migrations                  # migrations file directory
    ├── resources                   # additional file
    ├── services                    # services directory
    │   ├── data-processor          # app directory for data-processor service
    │   │  ├── constants            # Application global constants package
    │   │  ├── docs                 # swagger docs directory
    │   │  ├── driver               # Documentation files (alternatively `doc`)
    │   │  ├── packages             # Application packages
    │   │  ├── types                # Application global types package
    │   │  ├── utils                # Application utilities package
    │   │  └── main.go
    │   │
    │   └── dummy-proccess          # app directory for data-processor service
    │      ├── constants            # Application global constants package
    │      ├── driver               # Documentation files (alternatively `doc`)
    │      ├── packages             # Application packages
    │      ├── types                # Application global types package
    │      ├── utils                # Application utilities package
    │      └── main.go
    └── README.md

## How To Deploy
To deploy all services you need Docker already installed on your machine. You can install or learn more about Docker [here](https://www.docker.com/).

### Data Processor Service
For data-processor service, you can first build the image by executing following command:

    sudo docker build --tag data-processor .

After the image is successfully built, you can proceed to the next step to deploy service container by executing following command:

    sudo docker run -d --name data-processor -p 8081:8081 -e RABBITMQ_URL=amqp://lukas:lukas123@172.17.0.6:5672/ -e DB_URI=lukas:lukas123@tcp\(172.17.0.5:3306\)/simple_microservices?charset=utf8mb4\&parseTime=True\&loc=Local -e AUTH_JWT_PRIVATE_KEY=jwtRSA256.pem -e AUTH_JWT_PUBLIC_KEY=jwtRSA256.pub data-processor

#### Supported environment variable
Following is the list of supported environment variable for data-processor service

    APP_PORT                     # used to configure port used by http service. example value : 8080
    DB_URI                       # connection parameter to connect to mysql server, example value: user:password@tcp\(127.0.0.1:3306\)/simple_microservices?charset=utf8mb4\&parseTime=True\&loc=Local
    DB_DRIVER_TYPE               # Database driver type. You may leave this empty for it to use default value
    DB_MAX_OPEN_CONN             # Database connection pool configuration. You may leave this empty for it to use default value
    DB_MAX_CONN_LIFETIME         # Database connection pool configuration. You may leave this empty for it to use default value
    DB_MAX_IDLE_CONN             # Database connection pool configuration. You may leave this empty for it to use default value
    RABBITMQ_URL                 # connection parameter to connect to Rabbitmq, example value: amqp://user:password@localhost:5672/
    RABBITMQ_XCH_NAME            # Rabbitmq exchange name. You may leave this empty for it to use default value
    RABBITMQ_Q_NAME              # Rabbitmq queue name. You may leave this empty for it to use default value
    RABBITMQ_ROUTINGKEY          # Rabbitmq routing key. You may leave this empty for it to use default value
    RABBITMQ_CONSUMER_NAME       # Rabbitmq consumer namey. You may leave this empty for it to use default value
    AUTH_TYPE                    # Authentication configuration. You can set it to 2 possible values 'JWT' and 'STATIC'. the default value is 'JWT'
    AUTH_TOKENS                  # Use to set authorized token that can be use to access resources (only in STATIC authentication). Separate multiple value with semicolon ';'
    AUTH_JWT_PRIVATE_KEY         # Path to RSA private key (used if auth type is JWT)
    AUTH_JWT_PUBLIC_KEY          # Path to RSA public key (used if auth type is JWT)
    AUTH_JWT_EXPIRATION          # JWT token lifetime (expiration) (used if auth type is JWT). The default value is 24 hours. Example value : '10s' (10 seconds)


### Dummy Sensor
For dummy-sensor service the process is similar with data-processor service
To create image:

    sudo docker build --tag dummy-sensor .

To deploy container:

    sudo docker run -d --name dummy-sensor-humidity -p 8082:8082 -e APP_PORT=8082 -e DATA_SENSOR_TYPE=HUMIDITY -e RABBITMQ_URL=amqp://lukas:lukas123@172.17.0.6:5672/ dummy-sensor

#### Supported environment variable
Following is the list of supported environment variable for dummy-sensor service

    APP_PORT               # used to configure port used by http service. example value : 8080
    RABBITMQ_URL           # connection parameter to connect to Rabbitmq, example value : amqp://user:password@localhost:5672/
    RABBITMQ_XCH_NAME      # Rabbitmq exchange name. You may leave this empty for it to use default value
    RABBITMQ_Q_NAME        # Rabbitmq queue name. You may leave this empty for it to use default value
    RABBITMQ_ROUTINGKEY    # Rabbitmq routing key. You may leave this empty for it to use default value
    DATA_SENSOR_TYPE       # used to configure data sensor type. example value : TEMPERATURE
    DATA_SENSOR_ID1        # used to configure data ID1. example value : FO
    DATA_SENSOR_ID2        # used to configure data ID2. example value : 1




# Additional Resources

## How to generate RSA auth key pair
To generate RSA key pair we can use tool such as ssh-keygen

    ssh-keygen -t rsa -m PEM

    ssh-keygen -y -f yourkey.pem > key.pub