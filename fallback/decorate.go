package fallback

func DecorateRunnable[E error](
	fn func() error,
	fallback func(Context[any, any, E]) error,
	predicate FailurePredicate[any, any, E]) func() error {
	return func() error {
		ctx := execute[any, any](nil, func() (any, error) {
			return nil, fn()
		})
		if ok, failCtx := predicate(ctx); ok {
			return fallback(failCtx)
		}
		return ctx.returnError()
	}
}

func DecorateRunnableWithFailure[E error](
	fn func() error, fallback func(E) error) func() error {
	return DecorateRunnable(fn, func(ctx Context[any, any, E]) error {
		return fallback(ctx.Err())
	}, DefaultFailurePredicate[any, any, E]())
}

func DecorateRunnableByType[E error](
	fn func() error, fallback func() error) func() error {
	return DecorateRunnableWithFailure(fn, func(_ E) error { return fallback() })
}

func DecorateRunnableDefault(
	fn func() error, fallback func() error) func() error {
	return DecorateRunnableByType[error](fn, fallback)
}

func DecorateSupplier[R any, E error](
	fn func() (R, error),
	fallback func(Context[any, R, E]) (R, error),
	predicate FailurePredicate[any, R, E]) func() (R, error) {
	return func() (R, error) {
		ctx := execute[any, R](nil, func() (R, error) {
			return fn()
		})
		if ok, failCtx := predicate(ctx); ok {
			return fallback(failCtx)
		}
		return ctx.returnRetAndError()
	}
}

func DecorateSupplierWithFailure[R any, E error](
	fn func() (R, error), fallback func(R, E) (R, error)) func() (R, error) {
	return DecorateSupplier(fn, func(ctx Context[any, R, E]) (R, error) {
		return fallback(ctx.Ret(), ctx.Err())
	}, DefaultFailurePredicate[any, R, E]())
}

func DecorateSupplierByType[R any, E error](
	fn func() (R, error), fallback func() (R, error)) func() (R, error) {
	return DecorateSupplierWithFailure(fn, func(_ R, _ E) (R, error) { return fallback() })
}

func DecorateSupplierDefault[R any](
	fn func() (R, error), fallback func() (R, error)) func() (R, error) {
	return DecorateSupplierByType[R, error](fn, fallback)
}

func DecorateConsumer[T any, E error](
	fn func(T) error,
	fallback func(Context[T, any, E]) error,
	predicate FailurePredicate[T, any, E]) func(T) error {
	return func(t T) error {
		ctx := execute[T, any](t, func() (any, error) {
			return nil, fn(t)
		})
		if ok, failCtx := predicate(ctx); ok {
			return fallback(failCtx)
		}
		return ctx.returnError()
	}
}

func DecorateConsumerWithFailure[T any, E error](
	fn func(T) error, fallback func(T, E) error) func(T) error {
	return DecorateConsumer(fn, func(ctx Context[T, any, E]) error {
		return fallback(ctx.Param(), ctx.Err())
	}, DefaultFailurePredicate[T, any, E]())
}

func DecorateConsumerByType[T any, E error](
	fn func(T) error, fallback func(T) error) func(T) error {
	return DecorateConsumerWithFailure(fn, func(t T, _ E) error { return fallback(t) })
}

func DecorateConsumerDefault[T any](
	fn func(T) error, fallback func(T) error) func(T) error {
	return DecorateConsumerByType[T, error](fn, fallback)
}

func DecorateFunction[T any, R any, E error](
	fn func(T) (R, error),
	fallback func(Context[T, R, E]) (R, error),
	predicate FailurePredicate[T, R, E]) func(T) (R, error) {
	return func(t T) (R, error) {
		ctx := execute[T, R](t, func() (R, error) {
			return fn(t)
		})
		if ok, failCtx := predicate(ctx); ok {
			return fallback(failCtx)
		}
		return ctx.returnRetAndError()
	}
}

func DecorateFunctionWithFailure[T any, R any, E error](
	fn func(T) (R, error), fallback func(T, R, E) (R, error)) func(T) (R, error) {
	return DecorateFunction(fn, func(ctx Context[T, R, E]) (R, error) {
		return fallback(ctx.Param(), ctx.Ret(), ctx.Err())
	}, DefaultFailurePredicate[T, R, E]())
}

func DecorateFunctionByType[T any, R any, E error](
	fn func(T) (R, error), fallback func(T) (R, error)) func(T) (R, error) {
	return DecorateFunctionWithFailure(fn, func(t T, _ R, _ E) (R, error) { return fallback(t) })
}

func DecorateFunctionDefault[T any, R any](
	fn func(T) (R, error), fallback func(T) (R, error)) func(T) (R, error) {
	return DecorateFunctionByType[T, R, error](fn, fallback)
}
