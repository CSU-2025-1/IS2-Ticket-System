package rabbitmq

import "context"

type PublishQueueAdapter[T any] interface {
	Publish(msg T, opts ...PublishOption) error
	Close()
}

type Writer[T any] struct {
	adapter PublishQueueAdapter[T]
	option  PublishOption
}

func (w *Writer[T]) Send(ctx context.Context, msg T) error {
	return w.adapter.Publish(msg, w.option)
}

func (w *Writer[T]) Close() {
	w.adapter.Close()
}
