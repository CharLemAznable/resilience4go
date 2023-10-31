package timelimiter

import (
	"github.com/CharLemAznable/gofn/common"
	"github.com/CharLemAznable/gofn/consumer"
	"github.com/CharLemAznable/gofn/function"
	"github.com/CharLemAznable/gofn/runnable"
	"github.com/CharLemAznable/gofn/supplier"
)

func DecorateRunnable(limiter TimeLimiter, fn runnable.Runnable) runnable.Runnable {
	return func() error {
		_, err := limiter.execute(func() (any, error) {
			return nil, fn()
		})
		return err
	}
}

func DecorateSupplier[T any](limiter TimeLimiter, fn supplier.Supplier[T]) supplier.Supplier[T] {
	return func() (T, error) {
		ret, err := limiter.execute(func() (any, error) {
			return fn()
		})
		return common.CastQuietly[T](ret), err
	}
}

func DecorateConsumer[T any](limiter TimeLimiter, fn consumer.Consumer[T]) consumer.Consumer[T] {
	return func(t T) error {
		_, err := limiter.execute(func() (any, error) {
			return nil, fn(t)
		})
		return err
	}
}

func DecorateFunction[T any, R any](limiter TimeLimiter, fn function.Function[T, R]) function.Function[T, R] {
	return func(t T) (R, error) {
		ret, err := limiter.execute(func() (any, error) {
			return fn(t)
		})
		return common.CastQuietly[R](ret), err
	}
}
