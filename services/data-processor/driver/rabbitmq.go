package driver

import (
	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitmqOptions struct {
	URL          string
	ExchangeName string
	ExchangeType RabbitmqExchangeType
	QueueName    string
	RoutingKey   string
	ConsumerName string
}

type RabbitmqConnection struct {
	Options  *Options
	Instance *RabbitmqInstance
}

type RabbitmqInstance struct {
	Channel *amqp.Channel
	Queue   amqp.Queue
}

type RabbitmqExchangeType string

const (
	Direct  RabbitmqExchangeType = "direct"
	Topic   RabbitmqExchangeType = "topic"
	Headers RabbitmqExchangeType = "headers"
	Fanout  RabbitmqExchangeType = "fanout"
)

func (cp *RabbitmqConnection) Init() error {
	if rabbitConn, e := amqp.Dial(cp.Options.Rabbitmq.URL); e != nil {
		return e
	} else {
		if ch, e := rabbitConn.Channel(); e != nil {
			return e
		} else {
			if e := ch.ExchangeDeclare(
				cp.Options.Rabbitmq.ExchangeName,         // name
				string(cp.Options.Rabbitmq.ExchangeType), // type
				true,                                     // durable
				false,                                    // auto-deleted
				false,                                    // internal
				false,                                    // no-wait
				nil,                                      // arguments
			); e != nil {
				return e
			} else {
				if q, e := ch.QueueDeclare(
					cp.Options.Rabbitmq.QueueName, // name
					true,                          // durable
					false,                         // delete when unused
					false,                         // exclusive
					false,                         // no-wait
					nil,                           // arguments
				); e != nil {
					return e
				} else {
					ch.QueueBind(q.Name, cp.Options.Rabbitmq.RoutingKey, cp.Options.Rabbitmq.ExchangeName, false, nil)
					cp.Instance = &RabbitmqInstance{
						Channel: ch,
						Queue:   q,
					}
					return nil
				}
			}
		}
	}
}

func (cp *RabbitmqConnection) GetInstance() any {
	return cp.Instance
}

func (cp *RabbitmqConnection) GetWrapperInstance() any {
	return cp.Instance
}

func (cp *RabbitmqConnection) GetOptions() *Options {
	return cp.Options
}
