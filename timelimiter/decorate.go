package timelimiter

import "github.com/CharLemAznable/gogo/lang"

func DecorateRunnable(limiter TimeLimiter, fn func() error) func() error {
	return func() error {
		_, err := limiter.Execute(func() (any, error) {
			return nil, fn()
		})
		return err
	}
}

func DecorateSupplier[T any](limiter TimeLimiter, fn func() (T, error)) func() (T, error) {
	return func() (T, error) {
		ret, err := limiter.Execute(func() (any, error) {
			return fn()
		})
		return lang.CastQuietly[T](ret), err
	}
}

func DecorateConsumer[T any](limiter TimeLimiter, fn func(T) error) func(T) error {
	return func(t T) error {
		_, err := limiter.Execute(func() (any, error) {
			return nil, fn(t)
		})
		return err
	}
}

func DecorateFunction[T any, R any](limiter TimeLimiter, fn func(T) (R, error)) func(T) (R, error) {
	return func(t T) (R, error) {
		ret, err := limiter.Execute(func() (any, error) {
			return fn(t)
		})
		return lang.CastQuietly[R](ret), err
	}
}
