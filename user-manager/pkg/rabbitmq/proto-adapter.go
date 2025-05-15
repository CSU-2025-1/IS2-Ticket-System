package rabbitmq

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	amqp "github.com/rabbitmq/amqp091-go"
	"log/slog"
)

type ProtoQueueAdapter[T proto.Message] struct {
	conn        *amqp.Connection
	channel     *amqp.Channel
	queue       string
	constructor func() T
}

func NewProtoQueueAdapter[T proto.Message](cfg Config, constructor func() T) (*ProtoQueueAdapter[T], error) {
	conn, err := amqp.Dial(cfg.ToDSN())
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %s", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("failed to open a channel: %s", err)
	}

	return &ProtoQueueAdapter[T]{
		conn:        conn,
		channel:     ch,
		queue:       cfg.Queue,
		constructor: constructor,
	}, nil
}

func (p *ProtoQueueAdapter[T]) Consume(opts ...ConsumeOption) (<-chan T, error) {
	opt := DefaultConsumeOption
	if len(opts) > 0 {
		opt = opts[0]
	}

	msgs, err := p.channel.Consume(
		p.queue,
		opt.Consumer,
		opt.AutoAck,
		opt.Exclusive,
		opt.NoLocal,
		opt.NoWait,
		opt.Args,
	)
	if err != nil {
		return nil, err
	}

	channel := make(chan T)
	go func() {
		for d := range msgs {
			msg := p.constructor()
			if err := proto.Unmarshal(d.Body, msg); err != nil {
				slog.Error("[ Rabbit MQ failed unmarshal proto message ]", slog.String("err", err.Error()))
				continue
			}
			channel <- msg
			d.Ack(false)
		}
	}()

	return channel, nil
}

func (p *ProtoQueueAdapter[T]) Publish(msg T, opts ...PublishOption) error {
	opt := DefaultPublishOption
	if len(opts) > 0 {
		opt = opts[0]
	}

	data, err := proto.Marshal(msg)
	if err != nil {
		return fmt.Errorf("failed to marshal result: %s", err)
	}

	buf := make([]byte, len(data)+1)
	buf[0] = 0x01
	copy(buf[1:], data)

	return p.channel.Publish(
		opt.Exchange,
		p.queue,
		opt.Mandatory,
		opt.Immediate,
		amqp.Publishing{
			ContentType:  ContentTypeProtobuf,
			Body:         buf,
			DeliveryMode: amqp.Persistent,
		},
	)
}

func (p *ProtoQueueAdapter[T]) Close() {
	p.channel.Close()
	p.conn.Close()
}
