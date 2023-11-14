package fallback

import (
	"github.com/CharLemAznable/gofn/common"
	"github.com/CharLemAznable/gofn/consumer"
	"github.com/CharLemAznable/gofn/function"
	"github.com/CharLemAznable/gofn/runnable"
	"github.com/CharLemAznable/gofn/supplier"
)

func DecorateRunnable[E error](
	fn runnable.Runnable,
	fallback func(E) error,
	predicate FailurePredicate[E]) runnable.Runnable {
	return func() error {
		val := execute(func() (any, error) {
			return nil, fn()
		})
		if ok, err := predicate(val.err, val.panic); ok {
			return fallback(err)
		}
		return succeedReturn(val)
	}
}

func DecorateRunnableDefault[E error](
	fn runnable.Runnable,
	fallback func(E) error) runnable.Runnable {
	return DecorateRunnable(fn, fallback, defaultFailurePredicate[E]())
}

func DecorateSupplier[T any, E error](
	fn supplier.Supplier[T],
	fallback func(T, E) (T, error),
	predicate FailureResultPredicate[T, E]) supplier.Supplier[T] {
	return func() (T, error) {
		val := execute(func() (any, error) {
			return fn()
		})
		if ok, ret, err := predicate(common.CastQuietly[T](val.ret), val.err, val.panic); ok {
			return fallback(ret, err)
		}
		return succeedResultReturn[T](val)
	}
}

func DecorateSupplierDefault[T any, E error](
	fn supplier.Supplier[T],
	fallback func(T, E) (T, error)) supplier.Supplier[T] {
	return DecorateSupplier(fn, fallback, defaultFailureResultPredicate[T, E]())
}

func DecorateConsumer[T any, E error](
	fn consumer.Consumer[T],
	fallback func(T, E) error,
	predicate FailurePredicate[E]) consumer.Consumer[T] {
	return func(t T) error {
		val := execute(func() (any, error) {
			return nil, fn(t)
		})
		if ok, err := predicate(val.err, val.panic); ok {
			return fallback(t, err)
		}
		return succeedReturn(val)
	}
}

func DecorateConsumerDefault[T any, E error](
	fn consumer.Consumer[T],
	fallback func(T, E) error) consumer.Consumer[T] {
	return DecorateConsumer(fn, fallback, defaultFailurePredicate[E]())
}

func DecorateFunction[T any, R any, E error](
	fn function.Function[T, R],
	fallback func(T, R, E) (R, error),
	predicate FailureResultPredicate[R, E]) function.Function[T, R] {
	return func(t T) (R, error) {
		val := execute(func() (any, error) {
			return fn(t)
		})
		if ok, ret, err := predicate(common.CastQuietly[R](val.ret), val.err, val.panic); ok {
			return fallback(t, ret, err)
		}
		return succeedResultReturn[R](val)
	}
}

func DecorateFunctionDefault[T any, R any, E error](
	fn function.Function[T, R],
	fallback func(T, R, E) (R, error)) function.Function[T, R] {
	return DecorateFunction(fn, fallback, defaultFailureResultPredicate[R, E]())
}
