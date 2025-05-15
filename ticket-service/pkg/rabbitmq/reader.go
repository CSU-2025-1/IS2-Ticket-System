package rabbitmq

import (
	"context"
	"github.com/golang/protobuf/proto"
)

type ConsumeQueueAdapter[T any] interface {
	Consume(opts ...ConsumeOption) (<-chan T, error)
	Close()
}

type Reader[T any] struct {
	adapter ConsumeQueueAdapter[T]
	handler func(T)
	option  ConsumeOption
}

func CreateJsonReader[T any](cfg Config, handler func(T), opts ConsumeOption) (*Reader[T], error) {
	var reader Reader[T]

	var err error
	reader.adapter, err = NewJsonQueueAdapter[T](cfg)
	if err != nil {
		return nil, err
	}

	reader.handler = handler
	reader.option = opts

	return &reader, nil
}

func CreateProtoReader[T proto.Message](cfg Config, handler func(T), opts ConsumeOption, constructor func() T) (*Reader[T], error) {
	var reader Reader[T]

	var err error
	reader.adapter, err = NewProtoQueueAdapter[T](cfg, constructor)
	if err != nil {
		return nil, err
	}

	reader.handler = handler
	reader.option = opts

	return &reader, nil
}

func (r *Reader[T]) Run(ctx context.Context) error {
	channel, err := r.adapter.Consume(r.option)
	if err != nil {
		return err
	}

	for {
		select {
		case <-ctx.Done():
			return nil
		case msg, ok := <-channel:
			if !ok {
				return nil
			}
			r.handler(msg)
		}
	}
}

func (r *Reader[T]) Close() {
	r.adapter.Close()
}
