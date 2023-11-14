package fallback

import "github.com/CharLemAznable/gofn/common"

type channelValue struct {
	ret   any
	err   error
	panic any
}

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

type FailurePredicate[E error] func(err error, panic any) (bool, E)

// 默认失败断言: 仅判断error非空且类型匹配
func defaultFailurePredicate[E error]() FailurePredicate[E] {
	return func(err error, panic any) (bool, E) {
		if e, ok := panic.(E); ok {
			return ok, e
		}
		if e, ok := err.(E); ok {
			return ok, e
		}
		return false, common.Zero[E]()
	}
}

func succeedReturn(val *channelValue) error {
	if val.panic != nil {
		panic(val.panic)
	}
	return val.err
}

type FailureResultPredicate[T any, E error] func(ret T, err error, panic any) (bool, T, E)

// 带返回值的默认失败断言: 仅判断error非空且类型匹配
func defaultFailureResultPredicate[T any, E error]() FailureResultPredicate[T, E] {
	return func(ret T, err error, panic any) (bool, T, E) {
		okErr, e := defaultFailurePredicate[E]()(err, panic)
		return okErr, ret, e
	}
}

func succeedResultReturn[T any](val *channelValue) (T, error) {
	err := succeedReturn(val)
	return common.CastQuietly[T](val.ret), err
}
