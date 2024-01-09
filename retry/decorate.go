package retry

import (
	. "github.com/CharLemAznable/gogo/fn"
	"github.com/CharLemAznable/gogo/lang"
)

func DecorateRunnable(retry Retry, runnable Runnable) Runnable {
	return RunnableCast(func() error {
		_, err := retry.Execute(func() (any, error) {
			return nil, runnable.CheckedRun()
		})
		return err
	})
}

func DecorateRun(retry Retry, fn func()) func() {
	return DecorateRunnable(retry, RunnableOf(fn)).Run
}

func DecorateCheckedRun(retry Retry, fn func() error) func() error {
	return DecorateRunnable(retry, RunnableCast(fn)).CheckedRun
}

func DecorateSupplier[T any](retry Retry, supplier Supplier[T]) Supplier[T] {
	return SupplierCast(func() (T, error) {
		ret, err := retry.Execute(func() (any, error) {
			return supplier.CheckedGet()
		})
		return lang.CastQuietly[T](ret), err
	})
}

func DecorateGet[T any](retry Retry, fn func() T) func() T {
	return DecorateSupplier(retry, SupplierOf(fn)).Get
}

func DecorateCheckedGet[T any](retry Retry, fn func() (T, error)) func() (T, error) {
	return DecorateSupplier(retry, SupplierCast(fn)).CheckedGet
}

func DecorateConsumer[T any](retry Retry, consumer Consumer[T]) Consumer[T] {
	return ConsumerCast(func(t T) error {
		_, err := retry.Execute(func() (any, error) {
			return nil, consumer.CheckedAccept(t)
		})
		return err
	})
}

func DecorateAccept[T any](retry Retry, fn func(T)) func(T) {
	return DecorateConsumer(retry, ConsumerOf(fn)).Accept
}

func DecorateCheckedAccept[T any](retry Retry, fn func(T) error) func(T) error {
	return DecorateConsumer(retry, ConsumerCast(fn)).CheckedAccept
}

func DecorateFunction[T any, R any](retry Retry, function Function[T, R]) Function[T, R] {
	return FunctionCast(func(t T) (R, error) {
		ret, err := retry.Execute(func() (any, error) {
			return function.CheckedApply(t)
		})
		return lang.CastQuietly[R](ret), err
	})
}

func DecorateApply[T any, R any](retry Retry, fn func(T) R) func(T) R {
	return DecorateFunction(retry, FunctionOf(fn)).Apply
}

func DecorateCheckedApply[T any, R any](retry Retry, fn func(T) (R, error)) func(T) (R, error) {
	return DecorateFunction(retry, FunctionCast(fn)).CheckedApply
}
