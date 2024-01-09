package ratelimiter

import (
	. "github.com/CharLemAznable/gogo/fn"
	"github.com/CharLemAznable/gogo/lang"
)

func DecorateRunnable(limiter RateLimiter, runnable Runnable) Runnable {
	return RunnableCast(func() error {
		if err := limiter.AcquirePermission(); err != nil {
			return err
		}
		return runnable.CheckedRun()
	})
}

func DecorateRun(limiter RateLimiter, fn func()) func() {
	return DecorateRunnable(limiter, RunnableOf(fn)).Run
}

func DecorateCheckedRun(limiter RateLimiter, fn func() error) func() error {
	return DecorateRunnable(limiter, RunnableCast(fn)).CheckedRun
}

func DecorateSupplier[T any](limiter RateLimiter, supplier Supplier[T]) Supplier[T] {
	return SupplierCast(func() (T, error) {
		if err := limiter.AcquirePermission(); err != nil {
			return lang.Zero[T](), err
		}
		return supplier.CheckedGet()
	})
}

func DecorateGet[T any](limiter RateLimiter, fn func() T) func() T {
	return DecorateSupplier(limiter, SupplierOf(fn)).Get
}

func DecorateCheckedGet[T any](limiter RateLimiter, fn func() (T, error)) func() (T, error) {
	return DecorateSupplier(limiter, SupplierCast(fn)).CheckedGet
}

func DecorateConsumer[T any](limiter RateLimiter, consumer Consumer[T]) Consumer[T] {
	return ConsumerCast(func(t T) error {
		if err := limiter.AcquirePermission(); err != nil {
			return err
		}
		return consumer.CheckedAccept(t)
	})
}

func DecorateAccept[T any](limiter RateLimiter, fn func(T)) func(T) {
	return DecorateConsumer(limiter, ConsumerOf(fn)).Accept
}

func DecorateCheckedAccept[T any](limiter RateLimiter, fn func(T) error) func(T) error {
	return DecorateConsumer(limiter, ConsumerCast(fn)).CheckedAccept
}

func DecorateFunction[T any, R any](limiter RateLimiter, function Function[T, R]) Function[T, R] {
	return FunctionCast(func(t T) (R, error) {
		if err := limiter.AcquirePermission(); err != nil {
			return lang.Zero[R](), err
		}
		return function.CheckedApply(t)
	})
}

func DecorateApply[T any, R any](limiter RateLimiter, fn func(T) R) func(T) R {
	return DecorateFunction(limiter, FunctionOf(fn)).Apply
}

func DecorateCheckedApply[T any, R any](limiter RateLimiter, fn func(T) (R, error)) func(T) (R, error) {
	return DecorateFunction(limiter, FunctionCast(fn)).CheckedApply
}
