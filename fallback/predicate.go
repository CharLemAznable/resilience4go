package fallback

type FailurePredicate[T any, R any, E error] func(Context[T, R, error]) (bool, Context[T, R, E])

// DefaultFailurePredicate 默认失败断言: 仅判断error非空且类型匹配
func DefaultFailurePredicate[T any, R any, E error]() FailurePredicate[T, R, E] {
	return func(ctx Context[T, R, error]) (bool, Context[T, R, E]) {
		if e, ok := ctx.Err().(E); ok {
			return true, NewContext(ctx.Param(), ctx.Ret(), e, ctx.Panic())
		}
		return false, nil
	}
}
