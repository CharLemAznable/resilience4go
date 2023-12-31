package circuitbreaker

import "github.com/CharLemAznable/gogo/lang"

func DecorateRunnable(breaker CircuitBreaker, fn func() error) func() error {
	return func() error {
		_, err := breaker.Execute(func() (any, error) {
			return nil, fn()
		})
		return err
	}
}

func DecorateSupplier[T any](breaker CircuitBreaker, fn func() (T, error)) func() (T, error) {
	return func() (T, error) {
		ret, err := breaker.Execute(func() (any, error) {
			return fn()
		})
		return lang.CastQuietly[T](ret), err
	}
}

func DecorateConsumer[T any](breaker CircuitBreaker, fn func(T) error) func(T) error {
	return func(t T) error {
		_, err := breaker.Execute(func() (any, error) {
			return nil, fn(t)
		})
		return err
	}
}

func DecorateFunction[T any, R any](breaker CircuitBreaker, fn func(T) (R, error)) func(T) (R, error) {
	return func(t T) (R, error) {
		ret, err := breaker.Execute(func() (any, error) {
			return fn(t)
		})
		return lang.CastQuietly[R](ret), err
	}
}
