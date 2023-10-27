package timelimiter

import (
	"context"
	"fmt"
	"github.com/CharLemAznable/resilience4go/common"
)

type TimeLimiter interface {
	Name() string
	EventListener() EventListener

	execute(func() (any, error)) (any, error)
	onSuccess()
	onTimeout()
	onFailure(any)
}

func NewTimeLimiter(name string, configs ...ConfigBuilder) TimeLimiter {
	config := defaultConfig()
	for _, cfg := range configs {
		cfg(config)
	}
	return &timeLimiter{
		name:          name,
		config:        config,
		rootContext:   context.Background(),
		eventListener: newEventListener(),
	}
}

type timeLimiter struct {
	name          string
	config        *Config
	rootContext   context.Context
	eventListener EventListener
}

func (limiter *timeLimiter) Name() string {
	return limiter.name
}

func (limiter *timeLimiter) EventListener() EventListener {
	return limiter.eventListener
}

func (limiter *timeLimiter) execute(fn func() (any, error)) (any, error) {
	timeout, cancelFunc := context.WithTimeout(limiter.rootContext, limiter.config.timeoutDuration)
	defer cancelFunc()
	finished := make(chan *channelValue)
	panicked := make(common.Panicked)
	go func() {
		defer panicked.Recover()
		ret, err := fn()
		finished <- &channelValue{ret, err}
	}()
	select {
	case <-timeout.Done():
		limiter.onTimeout()
		return nil, &timeLimiterError{name: limiter.name}
	case result := <-finished:
		limiter.onSuccess()
		return result.ret, result.err
	case err := <-panicked.Caught():
		limiter.onFailure(err)
		return nil, common.PanicError(err)
	}
}

func (limiter *timeLimiter) onSuccess() {
	limiter.eventListener.consumeEvent(newSuccessEvent(limiter.name))
}

func (limiter *timeLimiter) onTimeout() {
	limiter.eventListener.consumeEvent(newTimeoutEvent(limiter.name))
}

func (limiter *timeLimiter) onFailure(error any) {
	limiter.eventListener.consumeEvent(newFailureEvent(limiter.name, error))
}

type timeLimiterError struct {
	name string
}

func (e *timeLimiterError) Error() string {
	return fmt.Sprintf("TimeLimiter '%s' recorded a timeout exception.", e.name)
}

type channelValue struct {
	ret interface{}
	err error
}
