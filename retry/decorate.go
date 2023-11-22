package retry

import (
	"github.com/CharLemAznable/gofn/common"
	"github.com/CharLemAznable/gofn/consumer"
	"github.com/CharLemAznable/gofn/function"
	"github.com/CharLemAznable/gofn/runnable"
	"github.com/CharLemAznable/gofn/supplier"
)

func DecorateRunnable(retry Retry, fn func() error) runnable.Runnable {
	return func() error {
		_, err := retry.execute(func() (any, error) {
			return nil, fn()
		})
		return err
	}
}

func DecorateSupplier[T any](retry Retry, fn func() (T, error)) supplier.Supplier[T] {
	return func() (T, error) {
		ret, err := retry.execute(func() (any, error) {
			return fn()
		})
		return common.CastQuietly[T](ret), err
	}
}

func DecorateConsumer[T any](retry Retry, fn func(T) error) consumer.Consumer[T] {
	return func(t T) error {
		_, err := retry.execute(func() (any, error) {
			return nil, fn(t)
		})
		return err
	}
}

func DecorateFunction[T any, R any](retry Retry, fn func(T) (R, error)) function.Function[T, R] {
	return func(t T) (R, error) {
		ret, err := retry.execute(func() (any, error) {
			return fn(t)
		})
		return common.CastQuietly[R](ret), err
	}
}
