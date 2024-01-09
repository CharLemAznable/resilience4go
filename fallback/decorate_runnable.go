package fallback

import . "github.com/CharLemAznable/gogo/fn"

func DecorateRunnable[E error](
	runnable Runnable,
	fallback func(Context[any, any, E]) error,
	predicate FailurePredicate[any, any, E]) Runnable {
	return RunnableCast(func() error {
		ctx := execute[any, any](nil, func() (any, error) {
			return nil, runnable.CheckedRun()
		})
		if ok, failCtx := predicate(ctx); ok {
			return fallback(failCtx)
		}
		return ctx.returnError()
	})
}

func DecorateRunnableWithFailure[E error](
	runnable Runnable, fallback func(E) error) Runnable {
	return DecorateRunnable(runnable, func(ctx Context[any, any, E]) error {
		return fallback(ctx.Err())
	}, DefaultFailurePredicate[any, any, E]())
}

func DecorateRunnableByType[E error](
	runnable Runnable, fallback func() error) Runnable {
	return DecorateRunnableWithFailure(runnable, func(_ E) error { return fallback() })
}

func DecorateRunnableDefault(
	runnable Runnable, fallback func() error) Runnable {
	return DecorateRunnableByType[error](runnable, fallback)
}

func DecorateCheckedRunWithFailure[E error](
	fn func() error, fallback func(E) error) func() error {
	return DecorateRunnableWithFailure(RunnableCast(fn), fallback).CheckedRun
}

func DecorateCheckedRunByType[E error](
	fn func() error, fallback func() error) func() error {
	return DecorateRunnableByType[E](RunnableCast(fn), fallback).CheckedRun
}

func DecorateCheckedRunDefault(
	fn func() error, fallback func() error) func() error {
	return DecorateRunnableDefault(RunnableCast(fn), fallback).CheckedRun
}
