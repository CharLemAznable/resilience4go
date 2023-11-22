package bulkhead

import (
	"github.com/CharLemAznable/gofn/common"
	"github.com/CharLemAznable/gofn/consumer"
	"github.com/CharLemAznable/gofn/function"
	"github.com/CharLemAznable/gofn/runnable"
	"github.com/CharLemAznable/gofn/supplier"
)

func DecorateRunnable(bulkhead Bulkhead, fn func() error) runnable.Runnable {
	return func() error {
		if err := bulkhead.acquire(); err != nil {
			return err
		}
		defer bulkhead.release()
		return fn()
	}
}

func DecorateSupplier[T any](bulkhead Bulkhead, fn func() (T, error)) supplier.Supplier[T] {
	return func() (T, error) {
		if err := bulkhead.acquire(); err != nil {
			return common.Zero[T](), err
		}
		defer bulkhead.release()
		return fn()
	}
}

func DecorateConsumer[T any](bulkhead Bulkhead, fn func(T) error) consumer.Consumer[T] {
	return func(t T) error {
		if err := bulkhead.acquire(); err != nil {
			return err
		}
		defer bulkhead.release()
		return fn(t)
	}
}

func DecorateFunction[T any, R any](bulkhead Bulkhead, fn func(T) (R, error)) function.Function[T, R] {
	return func(t T) (R, error) {
		if err := bulkhead.acquire(); err != nil {
			return common.Zero[R](), err
		}
		defer bulkhead.release()
		return fn(t)
	}
}
