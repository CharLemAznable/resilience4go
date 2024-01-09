package timelimiter

import (
	. "github.com/CharLemAznable/gogo/fn"
	"github.com/CharLemAznable/gogo/lang"
)

func DecorateRunnable(limiter TimeLimiter, runnable Runnable) Runnable {
	return RunnableCast(func() error {
		_, err := limiter.Execute(func() (any, error) {
			return nil, runnable.CheckedRun()
		})
		return err
	})
}

func DecorateRun(limiter TimeLimiter, fn func()) func() {
	return DecorateRunnable(limiter, RunnableOf(fn)).Run
}

func DecorateCheckedRun(limiter TimeLimiter, fn func() error) func() error {
	return DecorateRunnable(limiter, RunnableCast(fn)).CheckedRun
}

func DecorateSupplier[T any](limiter TimeLimiter, supplier Supplier[T]) Supplier[T] {
	return SupplierCast(func() (T, error) {
		ret, err := limiter.Execute(func() (any, error) {
			return supplier.CheckedGet()
		})
		return lang.CastQuietly[T](ret), err
	})
}

func DecorateGet[T any](limiter TimeLimiter, fn func() T) func() T {
	return DecorateSupplier(limiter, SupplierOf(fn)).Get
}

func DecorateCheckedGet[T any](limiter TimeLimiter, fn func() (T, error)) func() (T, error) {
	return DecorateSupplier(limiter, SupplierCast(fn)).CheckedGet
}

func DecorateConsumer[T any](limiter TimeLimiter, consumer Consumer[T]) Consumer[T] {
	return ConsumerCast(func(t T) error {
		_, err := limiter.Execute(func() (any, error) {
			return nil, consumer.CheckedAccept(t)
		})
		return err
	})
}

func DecorateAccept[T any](limiter TimeLimiter, fn func(T)) func(T) {
	return DecorateConsumer(limiter, ConsumerOf(fn)).Accept
}

func DecorateCheckedAccept[T any](limiter TimeLimiter, fn func(T) error) func(T) error {
	return DecorateConsumer(limiter, ConsumerCast(fn)).CheckedAccept
}

func DecorateFunction[T any, R any](limiter TimeLimiter, function Function[T, R]) Function[T, R] {
	return FunctionCast(func(t T) (R, error) {
		ret, err := limiter.Execute(func() (any, error) {
			return function.CheckedApply(t)
		})
		return lang.CastQuietly[R](ret), err
	})
}

func DecorateApply[T any, R any](limiter TimeLimiter, fn func(T) R) func(T) R {
	return DecorateFunction(limiter, FunctionOf(fn)).Apply
}

func DecorateCheckedApply[T any, R any](limiter TimeLimiter, fn func(T) (R, error)) func(T) (R, error) {
	return DecorateFunction(limiter, FunctionCast(fn)).CheckedApply
}
