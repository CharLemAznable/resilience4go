package fallback

import . "github.com/CharLemAznable/gogo/fn"

func DecorateSupplier[R any, E error](
	supplier Supplier[R],
	fallback func(Context[any, R, E]) (R, error),
	predicate FailurePredicate[any, R, E]) Supplier[R] {
	return SupplierCast(func() (R, error) {
		ctx := execute[any, R](nil, func() (R, error) {
			return supplier.CheckedGet()
		})
		if ok, failCtx := predicate(ctx); ok {
			return fallback(failCtx)
		}
		return ctx.returnRetAndError()
	})
}

func DecorateSupplierWithFailure[R any, E error](
	supplier Supplier[R], fallback func(R, E) (R, error)) Supplier[R] {
	return DecorateSupplier(supplier, func(ctx Context[any, R, E]) (R, error) {
		return fallback(ctx.Ret(), ctx.Err())
	}, DefaultFailurePredicate[any, R, E]())
}

func DecorateSupplierByType[R any, E error](
	supplier Supplier[R], fallback func() (R, error)) Supplier[R] {
	return DecorateSupplierWithFailure(supplier, func(_ R, _ E) (R, error) { return fallback() })
}

func DecorateSupplierDefault[R any](
	supplier Supplier[R], fallback func() (R, error)) Supplier[R] {
	return DecorateSupplierByType[R, error](supplier, fallback)
}

func DecorateCheckedGetWithFailure[R any, E error](
	fn func() (R, error), fallback func(R, E) (R, error)) func() (R, error) {
	return DecorateSupplierWithFailure(SupplierCast(fn), fallback).CheckedGet
}

func DecorateCheckedGetByType[R any, E error](
	fn func() (R, error), fallback func() (R, error)) func() (R, error) {
	return DecorateSupplierByType[R, E](SupplierCast(fn), fallback).CheckedGet
}

func DecorateCheckedGetDefault[R any](
	fn func() (R, error), fallback func() (R, error)) func() (R, error) {
	return DecorateSupplierDefault(SupplierCast(fn), fallback).CheckedGet
}
