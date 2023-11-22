package cache

import (
	"github.com/CharLemAznable/gofn/function"
)

func DecorateFunction[T any, R any](cache Cache[T, R], fn func(T) (R, error)) function.Function[T, R] {
	return func(t T) (R, error) {
		return cache.getOrLoad(t, func() (R, error) { return fn(t) })
	}
}
