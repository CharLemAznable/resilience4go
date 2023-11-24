package cache

func DecorateFunction[T any, R any](cache Cache[T, R], fn func(T) (R, error)) func(T) (R, error) {
	return func(t T) (R, error) { return cache.GetOrLoad(t, fn) }
}
