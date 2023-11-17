package fallback

type Context[T any, R any, E error] interface {
	Param() T
	Ret() R
	Err() E
	Panic() any

	returnError() E
	returnRetAndError() (R, E)
}

func NewContext[T any, R any, E error](param T, ret R, err E, panic any) Context[T, R, E] {
	return &context[T, R, E]{param: param, ret: ret, err: err, panic: panic}
}

type context[T any, R any, E error] struct {
	param T
	ret   R
	err   E
	panic any
}

func (c *context[T, R, E]) Param() T {
	return c.param
}

func (c *context[T, R, E]) Ret() R {
	return c.ret
}

func (c *context[T, R, E]) Err() E {
	return c.err
}

func (c *context[T, R, E]) Panic() any {
	return c.panic
}

func (c *context[T, R, E]) returnError() E {
	if c.panic != nil {
		panic(c.panic)
	}
	return c.err
}

func (c *context[T, R, E]) returnRetAndError() (R, E) {
	if c.panic != nil {
		panic(c.panic)
	}
	return c.ret, c.err
}
