package cache

import . "github.com/CharLemAznable/gogo/fn"

func DecorateFunction[T any, R any](cache Cache[T, R], function Function[T, R]) Function[T, R] {
	return FunctionCast(func(t T) (R, error) { return cache.GetOrLoad(t, function.CheckedApply) })
}

func DecorateApply[T any, R any](cache Cache[T, R], fn func(T) R) func(T) R {
	return DecorateFunction(cache, FunctionOf(fn)).Apply
}

func DecorateCheckedApply[T any, R any](cache Cache[T, R], fn func(T) (R, error)) func(T) (R, error) {
	return DecorateFunction(cache, FunctionCast(fn)).CheckedApply
}
