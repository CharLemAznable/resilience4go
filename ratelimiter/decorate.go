package ratelimiter

import (
	"github.com/CharLemAznable/gofn/common"
	"github.com/CharLemAznable/gofn/consumer"
	"github.com/CharLemAznable/gofn/function"
	"github.com/CharLemAznable/gofn/runnable"
	"github.com/CharLemAznable/gofn/supplier"
)

func DecorateRunnable(limiter RateLimiter, fn func() error) runnable.Runnable {
	return func() error {
		if err := limiter.acquirePermission(); err != nil {
			return err
		}
		return fn()
	}
}

func DecorateSupplier[T any](limiter RateLimiter, fn func() (T, error)) supplier.Supplier[T] {
	return func() (T, error) {
		if err := limiter.acquirePermission(); err != nil {
			return common.Zero[T](), err
		}
		return fn()
	}
}

func DecorateConsumer[T any](limiter RateLimiter, fn func(T) error) consumer.Consumer[T] {
	return func(t T) error {
		if err := limiter.acquirePermission(); err != nil {
			return err
		}
		return fn(t)
	}
}

func DecorateFunction[T any, R any](limiter RateLimiter, fn func(T) (R, error)) function.Function[T, R] {
	return func(t T) (R, error) {
		if err := limiter.acquirePermission(); err != nil {
			return common.Zero[R](), err
		}
		return fn(t)
	}
}
