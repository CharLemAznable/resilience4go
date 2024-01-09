package circuitbreaker

import (
	. "github.com/CharLemAznable/gogo/fn"
	"github.com/CharLemAznable/gogo/lang"
)

func DecorateRunnable(breaker CircuitBreaker, runnable Runnable) Runnable {
	return RunnableCast(func() error {
		_, err := breaker.Execute(func() (any, error) {
			return nil, runnable.CheckedRun()
		})
		return err
	})
}

func DecorateRun(breaker CircuitBreaker, fn func()) func() {
	return DecorateRunnable(breaker, RunnableOf(fn)).Run
}

func DecorateCheckedRun(breaker CircuitBreaker, fn func() error) func() error {
	return DecorateRunnable(breaker, RunnableCast(fn)).CheckedRun
}

func DecorateSupplier[T any](breaker CircuitBreaker, supplier Supplier[T]) Supplier[T] {
	return SupplierCast(func() (T, error) {
		ret, err := breaker.Execute(func() (any, error) {
			return supplier.CheckedGet()
		})
		return lang.CastQuietly[T](ret), err
	})
}

func DecorateGet[T any](breaker CircuitBreaker, fn func() T) func() T {
	return DecorateSupplier(breaker, SupplierOf(fn)).Get
}

func DecorateCheckedGet[T any](breaker CircuitBreaker, fn func() (T, error)) func() (T, error) {
	return DecorateSupplier(breaker, SupplierCast(fn)).CheckedGet
}

func DecorateConsumer[T any](breaker CircuitBreaker, consumer Consumer[T]) Consumer[T] {
	return ConsumerCast(func(t T) error {
		_, err := breaker.Execute(func() (any, error) {
			return nil, consumer.CheckedAccept(t)
		})
		return err
	})
}

func DecorateAccept[T any](breaker CircuitBreaker, fn func(T)) func(T) {
	return DecorateConsumer(breaker, ConsumerOf(fn)).Accept
}

func DecorateCheckedAccept[T any](breaker CircuitBreaker, fn func(T) error) func(T) error {
	return DecorateConsumer(breaker, ConsumerCast(fn)).CheckedAccept
}

func DecorateFunction[T any, R any](breaker CircuitBreaker, function Function[T, R]) Function[T, R] {
	return FunctionCast(func(t T) (R, error) {
		ret, err := breaker.Execute(func() (any, error) {
			return function.CheckedApply(t)
		})
		return lang.CastQuietly[R](ret), err
	})
}

func DecorateApply[T any, R any](breaker CircuitBreaker, fn func(T) R) func(T) R {
	return DecorateFunction(breaker, FunctionOf(fn)).Apply
}

func DecorateCheckedApply[T any, R any](breaker CircuitBreaker, fn func(T) (R, error)) func(T) (R, error) {
	return DecorateFunction(breaker, FunctionCast(fn)).CheckedApply
}
