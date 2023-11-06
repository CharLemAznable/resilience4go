package fallback

import (
	"github.com/CharLemAznable/gofn/consumer"
	"github.com/CharLemAznable/gofn/function"
	"github.com/CharLemAznable/gofn/runnable"
	"github.com/CharLemAznable/gofn/supplier"
)

func DecorateRunnable[E error](fn runnable.Runnable, fallback func(E) error) runnable.Runnable {
	return func() error {
		val := execute(func() (any, error) {
			return nil, fn()
		})
		if e, ok := castErr[E](val); ok {
			return fallback(e)
		}
		_, err := result[any](val)
		return err
	}
}

func DecorateSupplier[T any, E error](fn supplier.Supplier[T], fallback func(E) (T, error)) supplier.Supplier[T] {
	return func() (T, error) {
		val := execute(func() (any, error) {
			return fn()
		})
		if e, ok := castErr[E](val); ok {
			return fallback(e)
		}
		return result[T](val)
	}
}

func DecorateConsumer[T any, E error](fn consumer.Consumer[T], fallback func(E) error) consumer.Consumer[T] {
	return func(t T) error {
		val := execute(func() (any, error) {
			return nil, fn(t)
		})
		if e, ok := castErr[E](val); ok {
			return fallback(e)
		}
		_, err := result[any](val)
		return err
	}
}

func DecorateFunction[T any, R any, E error](fn function.Function[T, R], fallback func(E) (R, error)) function.Function[T, R] {
	return func(t T) (R, error) {
		val := execute(func() (any, error) {
			return fn(t)
		})
		if e, ok := castErr[E](val); ok {
			return fallback(e)
		}
		return result[R](val)
	}
}
