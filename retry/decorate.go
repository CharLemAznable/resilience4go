package retry

import "github.com/CharLemAznable/gogo/lang"

func DecorateRunnable(retry Retry, fn func() error) func() error {
	return func() error {
		_, err := retry.Execute(func() (any, error) {
			return nil, fn()
		})
		return err
	}
}

func DecorateSupplier[T any](retry Retry, fn func() (T, error)) func() (T, error) {
	return func() (T, error) {
		ret, err := retry.Execute(func() (any, error) {
			return fn()
		})
		return lang.CastQuietly[T](ret), err
	}
}

func DecorateConsumer[T any](retry Retry, fn func(T) error) func(T) error {
	return func(t T) error {
		_, err := retry.Execute(func() (any, error) {
			return nil, fn(t)
		})
		return err
	}
}

func DecorateFunction[T any, R any](retry Retry, fn func(T) (R, error)) func(T) (R, error) {
	return func(t T) (R, error) {
		ret, err := retry.Execute(func() (any, error) {
			return fn(t)
		})
		return lang.CastQuietly[R](ret), err
	}
}
