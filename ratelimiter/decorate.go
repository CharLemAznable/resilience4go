package ratelimiter

import "github.com/CharLemAznable/gogo/lang"

func DecorateRunnable(limiter RateLimiter, fn func() error) func() error {
	return func() error {
		if err := limiter.AcquirePermission(); err != nil {
			return err
		}
		return fn()
	}
}

func DecorateSupplier[T any](limiter RateLimiter, fn func() (T, error)) func() (T, error) {
	return func() (T, error) {
		if err := limiter.AcquirePermission(); err != nil {
			return lang.Zero[T](), err
		}
		return fn()
	}
}

func DecorateConsumer[T any](limiter RateLimiter, fn func(T) error) func(T) error {
	return func(t T) error {
		if err := limiter.AcquirePermission(); err != nil {
			return err
		}
		return fn(t)
	}
}

func DecorateFunction[T any, R any](limiter RateLimiter, fn func(T) (R, error)) func(T) (R, error) {
	return func(t T) (R, error) {
		if err := limiter.AcquirePermission(); err != nil {
			return lang.Zero[R](), err
		}
		return fn(t)
	}
}
