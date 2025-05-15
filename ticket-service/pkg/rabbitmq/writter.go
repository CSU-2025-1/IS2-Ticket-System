package rabbitmq

import (
	"context"
	"github.com/golang/protobuf/proto"
)

type PublishQueueAdapter[T any] interface {
	Publish(msg T, opts ...PublishOption) error
	Close()
}

type Writer[T any] struct {
	adapter PublishQueueAdapter[T]
	option  PublishOption
}

func CreateJsonWriter[T any](cfg Config, option PublishOption) (*Writer[T], error) {
	var writer Writer[T]

	var err error
	writer.adapter, err = NewJsonQueueAdapter[T](cfg)
	if err != nil {
		return nil, err
	}

	writer.option = option

	return &writer, nil
}

func CreateProtoWriter[T proto.Message](cfg Config, option PublishOption) (*Writer[T], error) {
	var writer Writer[T]

	var err error
	writer.adapter, err = NewProtoQueueAdapter[T](cfg, func() T {
		var zero T
		return zero
	})
	if err != nil {
		return nil, err
	}

	writer.option = option

	return &writer, nil
}

func (w *Writer[T]) Send(ctx context.Context, msg T) error {
	return w.adapter.Publish(msg, w.option)
}

func (w *Writer[T]) Close() {
	w.adapter.Close()
}
