package fallback

import . "github.com/CharLemAznable/gogo/fn"

func DecorateFunction[T any, R any, E error](
	function Function[T, R],
	fallback func(Context[T, R, E]) (R, error),
	predicate FailurePredicate[T, R, E]) Function[T, R] {
	return FunctionCast(func(t T) (R, error) {
		ctx := execute[T, R](t, func() (R, error) {
			return function.CheckedApply(t)
		})
		if ok, failCtx := predicate(ctx); ok {
			return fallback(failCtx)
		}
		return ctx.returnRetAndError()
	})
}

func DecorateFunctionWithFailure[T any, R any, E error](
	function Function[T, R], fallback func(T, R, E) (R, error)) Function[T, R] {
	return DecorateFunction(function, func(ctx Context[T, R, E]) (R, error) {
		return fallback(ctx.Param(), ctx.Ret(), ctx.Err())
	}, DefaultFailurePredicate[T, R, E]())
}

func DecorateFunctionByType[T any, R any, E error](
	function Function[T, R], fallback func(T) (R, error)) Function[T, R] {
	return DecorateFunctionWithFailure(function, func(t T, _ R, _ E) (R, error) { return fallback(t) })
}

func DecorateFunctionDefault[T any, R any](
	function Function[T, R], fallback func(T) (R, error)) Function[T, R] {
	return DecorateFunctionByType[T, R, error](function, fallback)
}

func DecorateCheckedApplyWithFailure[T any, R any, E error](
	fn func(T) (R, error), fallback func(T, R, E) (R, error)) func(T) (R, error) {
	return DecorateFunctionWithFailure(FunctionCast(fn), fallback).CheckedApply
}

func DecorateCheckedApplyByType[T any, R any, E error](
	fn func(T) (R, error), fallback func(T) (R, error)) func(T) (R, error) {
	return DecorateFunctionByType[T, R, E](FunctionCast(fn), fallback).CheckedApply
}

func DecorateCheckedApplyDefault[T any, R any](
	fn func(T) (R, error), fallback func(T) (R, error)) func(T) (R, error) {
	return DecorateFunctionDefault(FunctionCast(fn), fallback).CheckedApply
}
