package fallback

import (
	"github.com/CharLemAznable/gofn/consumer"
	"github.com/CharLemAznable/gofn/function"
	"github.com/CharLemAznable/gofn/runnable"
	"github.com/CharLemAznable/gofn/supplier"
)

func DecorateRunnable[E error](fn runnable.Runnable, fallback func(E) error) runnable.Runnable {
	return func() error {
		err := fn()
		if e, ok := err.(E); ok {
			return fallback(e)
		}
		return err
	}
}

func DecorateSupplier[T any, E error](fn supplier.Supplier[T], fallback func(E) (T, error)) supplier.Supplier[T] {
	return func() (T, error) {
		t, err := fn()
		if e, ok := err.(E); ok {
			return fallback(e)
		}
		return t, err
	}
}

func DecorateConsumer[T any, E error](fn consumer.Consumer[T], fallback func(E) error) consumer.Consumer[T] {
	return func(t T) error {
		err := fn(t)
		if e, ok := err.(E); ok {
			return fallback(e)
		}
		return err
	}
}

func DecorateFunction[T any, R any, E error](fn function.Function[T, R], fallback func(E) (R, error)) function.Function[T, R] {
	return func(t T) (R, error) {
		r, err := fn(t)
		if e, ok := err.(E); ok {
			return fallback(e)
		}
		return r, err
	}
}
