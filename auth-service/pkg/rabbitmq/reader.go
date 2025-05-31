package rabbitmq

import "context"

type ConsumeQueueAdapter[T any] interface {
	Consume(opts ...ConsumeOption) (<-chan T, error)
	Close()
}

type Reader[T any] struct {
	adapter ConsumeQueueAdapter[T]
	handler func(T)
	option  ConsumeOption
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
			go r.handler(msg)
		}
	}
}

func (r *Reader[T]) Close() {
	r.adapter.Close()
}
