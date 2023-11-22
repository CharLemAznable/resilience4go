package ratelimiter

import "github.com/CharLemAznable/gofn/common"

func DecorateRunnable(limiter RateLimiter, fn func() error) func() error {
	return func() error {
		if err := limiter.acquirePermission(); err != nil {
			return err
		}
		return fn()
	}
}

func DecorateSupplier[T any](limiter RateLimiter, fn func() (T, error)) func() (T, error) {
	return func() (T, error) {
		if err := limiter.acquirePermission(); err != nil {
			return common.Zero[T](), err
		}
		return fn()
	}
}

func DecorateConsumer[T any](limiter RateLimiter, fn func(T) error) func(T) error {
	return func(t T) error {
		if err := limiter.acquirePermission(); err != nil {
			return err
		}
		return fn(t)
	}
}

func DecorateFunction[T any, R any](limiter RateLimiter, fn func(T) (R, error)) func(T) (R, error) {
	return func(t T) (R, error) {
		if err := limiter.acquirePermission(); err != nil {
			return common.Zero[R](), err
		}
		return fn(t)
	}
}
