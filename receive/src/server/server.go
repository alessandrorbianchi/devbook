package server

import (
	"github.com/streadway/amqp"
)

func ConectarNoRabbit(appRabbitAddr string) (*amqp.Connection, error) {
	conn, err := amqp.Dial(appRabbitAddr)
	if err != nil {
		return nil, err
	}

	return conn, nil
}
