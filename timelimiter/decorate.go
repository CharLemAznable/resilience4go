package timelimiter

import "github.com/CharLemAznable/gofn/common"

func DecorateRunnable(limiter TimeLimiter, fn func() error) func() error {
	return func() error {
		_, err := limiter.execute(func() (any, error) {
			return nil, fn()
		})
		return err
	}
}

func DecorateSupplier[T any](limiter TimeLimiter, fn func() (T, error)) func() (T, error) {
	return func() (T, error) {
		ret, err := limiter.execute(func() (any, error) {
			return fn()
		})
		return common.CastQuietly[T](ret), err
	}
}

func DecorateConsumer[T any](limiter TimeLimiter, fn func(T) error) func(T) error {
	return func(t T) error {
		_, err := limiter.execute(func() (any, error) {
			return nil, fn(t)
		})
		return err
	}
}

func DecorateFunction[T any, R any](limiter TimeLimiter, fn func(T) (R, error)) func(T) (R, error) {
	return func(t T) (R, error) {
		ret, err := limiter.execute(func() (any, error) {
			return fn(t)
		})
		return common.CastQuietly[R](ret), err
	}
}
