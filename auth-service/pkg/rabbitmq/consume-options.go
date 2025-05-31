package rabbitmq

import amqp "github.com/rabbitmq/amqp091-go"

type ConsumeOption struct {
	Consumer  string
	AutoAck   bool
	Exclusive bool
	NoLocal   bool
	NoWait    bool
	Args      amqp.Table
}

var DefaultConsumeOption = ConsumeOption{
	Consumer:  "",
	AutoAck:   false,
	Exclusive: false,
	NoLocal:   false,
	NoWait:    false,
	Args:      nil,
}

func (c *ConsumeOption) Option() (string, bool, bool, bool, bool, amqp.Table) {
	return c.Consumer, c.AutoAck, c.Exclusive, c.NoLocal, c.NoWait, c.Args
}
