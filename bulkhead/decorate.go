package bulkhead

import "github.com/CharLemAznable/ge"

func DecorateRunnable(bulkhead Bulkhead, fn func() error) func() error {
	return func() error {
		if err := bulkhead.Acquire(); err != nil {
			return err
		}
		defer bulkhead.Release()
		return fn()
	}
}

func DecorateSupplier[T any](bulkhead Bulkhead, fn func() (T, error)) func() (T, error) {
	return func() (T, error) {
		if err := bulkhead.Acquire(); err != nil {
			return ge.Zero[T](), err
		}
		defer bulkhead.Release()
		return fn()
	}
}

func DecorateConsumer[T any](bulkhead Bulkhead, fn func(T) error) func(T) error {
	return func(t T) error {
		if err := bulkhead.Acquire(); err != nil {
			return err
		}
		defer bulkhead.Release()
		return fn(t)
	}
}

func DecorateFunction[T any, R any](bulkhead Bulkhead, fn func(T) (R, error)) func(T) (R, error) {
	return func(t T) (R, error) {
		if err := bulkhead.Acquire(); err != nil {
			return ge.Zero[R](), err
		}
		defer bulkhead.Release()
		return fn(t)
	}
}
