package fallback

import "github.com/CharLemAznable/gogo/lang"

func execute[T any, R any](param T, fn func() (R, error)) Context[T, R, error] {
	ctx := &context[T, R, error]{param: param}
	finished := make(chan *context[T, R, error])
	panicked := make(lang.Panicked)
	go func() {
		defer panicked.Recover()
		ctx.ret, ctx.err = fn()
		finished <- ctx
	}()
	select {
	case result := <-finished:
		return result
	case ctx.panic = <-panicked.Caught():
		return ctx
	}
}
