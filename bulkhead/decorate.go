package bulkhead

import (
	. "github.com/CharLemAznable/gogo/fn"
	"github.com/CharLemAznable/gogo/lang"
)

func DecorateRunnable(bulkhead Bulkhead, runnable Runnable) Runnable {
	return RunnableCast(func() error {
		if err := bulkhead.Acquire(); err != nil {
			return err
		}
		defer bulkhead.Release()
		return runnable.CheckedRun()
	})
}

func DecorateRun(bulkhead Bulkhead, fn func()) func() {
	return DecorateRunnable(bulkhead, RunnableOf(fn)).Run
}

func DecorateCheckedRun(bulkhead Bulkhead, fn func() error) func() error {
	return DecorateRunnable(bulkhead, RunnableCast(fn)).CheckedRun
}

func DecorateSupplier[T any](bulkhead Bulkhead, supplier Supplier[T]) Supplier[T] {
	return SupplierCast(func() (T, error) {
		if err := bulkhead.Acquire(); err != nil {
			return lang.Zero[T](), err
		}
		defer bulkhead.Release()
		return supplier.CheckedGet()
	})
}

func DecorateGet[T any](bulkhead Bulkhead, fn func() T) func() T {
	return DecorateSupplier(bulkhead, SupplierOf(fn)).Get
}

func DecorateCheckedGet[T any](bulkhead Bulkhead, fn func() (T, error)) func() (T, error) {
	return DecorateSupplier(bulkhead, SupplierCast(fn)).CheckedGet
}

func DecorateConsumer[T any](bulkhead Bulkhead, consumer Consumer[T]) Consumer[T] {
	return ConsumerCast(func(t T) error {
		if err := bulkhead.Acquire(); err != nil {
			return err
		}
		defer bulkhead.Release()
		return consumer.CheckedAccept(t)
	})
}

func DecorateAccept[T any](bulkhead Bulkhead, fn func(T)) func(T) {
	return DecorateConsumer(bulkhead, ConsumerOf(fn)).Accept
}

func DecorateCheckedAccept[T any](bulkhead Bulkhead, fn func(T) error) func(T) error {
	return DecorateConsumer(bulkhead, ConsumerCast(fn)).CheckedAccept
}

func DecorateFunction[T any, R any](bulkhead Bulkhead, function Function[T, R]) Function[T, R] {
	return FunctionCast(func(t T) (R, error) {
		if err := bulkhead.Acquire(); err != nil {
			return lang.Zero[R](), err
		}
		defer bulkhead.Release()
		return function.CheckedApply(t)
	})
}

func DecorateApply[T any, R any](bulkhead Bulkhead, fn func(T) R) func(T) R {
	return DecorateFunction(bulkhead, FunctionOf(fn)).Apply
}

func DecorateCheckedApply[T any, R any](bulkhead Bulkhead, fn func(T) (R, error)) func(T) (R, error) {
	return DecorateFunction(bulkhead, FunctionCast(fn)).CheckedApply
}
