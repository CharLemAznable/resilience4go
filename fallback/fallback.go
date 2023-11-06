package fallback

import "github.com/CharLemAznable/gofn/common"

func execute(fn func() (any, error)) *channelValue {
	finished := make(chan *channelValue)
	panicked := make(common.Panicked)
	go func() {
		defer panicked.Recover()
		ret, err := fn()
		finished <- &channelValue{ret, err, nil}
	}()
	select {
	case result := <-finished:
		return result
	case err := <-panicked.Caught():
		return &channelValue{nil, nil, err}
	}
}

func castErr[E error](val *channelValue) (E, bool) {
	if e, ok := val.panic.(E); ok {
		return e, ok
	}
	if e, ok := val.err.(E); ok {
		return e, ok
	}
	return common.Zero[E](), false
}

func result[T any](val *channelValue) (T, error) {
	if val.panic != nil {
		panic(val.panic)
	}
	return common.CastQuietly[T](val.ret), val.err
}

type channelValue struct {
	ret   any
	err   error
	panic any
}
