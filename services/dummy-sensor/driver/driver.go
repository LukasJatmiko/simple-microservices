package driver

type DriverType string

type Options struct {
	Type     DriverType
	Gorm     GormOptions
	Rabbitmq RabbitmqOptions
}

type Driver interface {
	Init() error
	GetInstance() any
	GetWrapperInstance() any
	GetOptions() *Options
}

const (
	GORM     DriverType = "GORM"
	GRPC     DriverType = "GRPC"
	RABBITMQ DriverType = "RABBITMQ"
)

func NewDriver(opts *Options) Driver {
	switch opts.Type {

	//if mysql
	case RABBITMQ:
		return &RabbitmqConnection{
			Options: opts,
		}

	//set default driver to postgres
	default:
		return &GormConnection{
			Options: opts,
		}
	}
}
