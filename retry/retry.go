package retry

import (
	"fmt"
	"github.com/CharLemAznable/gofn/common"
	"time"
)

type Retry interface {
	Name() string
	EventListener() EventListener

	execute(func() (any, error)) (any, error)
}

func NewRetry(name string, configs ...ConfigBuilder) Retry {
	config := defaultConfig()
	for _, cfg := range configs {
		cfg(config)
	}
	return &retry{
		name:          name,
		config:        config,
		eventListener: newEventListener(),
	}
}

type retry struct {
	name          string
	config        *Config
	eventListener EventListener
}

func (r *retry) Name() string {
	return r.name
}

func (r *retry) EventListener() EventListener {
	return r.eventListener
}

func (r *retry) execute(fn func() (any, error)) (any, error) {
	context := r.executeOnce(fn)
	if r.testResult(context) {
		return r.returnResult(context)
	}
	numOfAttempts := 1
	for ; numOfAttempts < r.config.maxAttempts; numOfAttempts++ {
		waitDuration := r.config.waitIntervalFunctionFn(numOfAttempts)
		r.publishEvent(newRetryEvent(r.name,
			numOfAttempts, context.ret, context.err, waitDuration))
		time.Sleep(waitDuration)

		context = r.executeOnce(fn)
		if r.testResult(context) {
			r.publishEvent(newSuccessEvent(r.name,
				numOfAttempts, context.ret, context.err))
			return r.returnResult(context)
		}
	}
	r.publishEvent(newErrorEvent(r.name,
		numOfAttempts, context.ret, context.err))
	if r.config.failAfterMaxAttempts {
		context.err = common.DefaultErrorFn(context.err, func() error {
			return &MaxRetriesExceeded{name: r.name, maxAttempts: r.config.maxAttempts}
		})
	}
	return r.returnResult(context)
}

func (r *retry) executeOnce(fn func() (any, error)) *channelValue {
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
		return &channelValue{nil, common.WrapPanic(err), err}
	}
}

func (r *retry) testResult(result *channelValue) bool {
	return !r.config.recordResultPredicateFn(result.ret, result.err)
}

func (r *retry) returnResult(result *channelValue) (any, error) {
	if result.panic != nil {
		panic(result.panic)
	}
	return result.ret, result.err
}

func (r *retry) publishEvent(event Event) {
	r.eventListener.consumeEvent(event)
}

type channelValue struct {
	ret   any
	err   error
	panic any
}

type MaxRetriesExceeded struct {
	name        string
	maxAttempts int
}

func (e *MaxRetriesExceeded) Error() string {
	return fmt.Sprintf("Retry '%s' has exhausted all attempts (%d)", e.name, e.maxAttempts)
}
