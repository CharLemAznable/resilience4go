package timelimiter

import (
	"context"
	"fmt"
	"github.com/CharLemAznable/gogo/lang"
)

type TimeLimiter interface {
	Name() string
	Metrics() Metrics
	EventListener() EventListener
	Execute(func() (any, error)) (any, error)
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
		metrics:       newMetrics(),
		eventListener: newEventListener(),
	}
}

type timeLimiter struct {
	name          string
	config        *Config
	rootContext   context.Context
	metrics       *metrics
	eventListener *eventListener
}

func (limiter *timeLimiter) Name() string {
	return limiter.name
}

func (limiter *timeLimiter) Metrics() Metrics {
	return limiter.metrics
}

func (limiter *timeLimiter) EventListener() EventListener {
	return limiter.eventListener
}

func (limiter *timeLimiter) Execute(fn func() (any, error)) (any, error) {
	timeout, cancelFunc := context.WithTimeout(limiter.rootContext, limiter.config.timeoutDuration)
	defer cancelFunc()
	finished := make(chan *channelValue)
	panicked := make(lang.Panicked)
	go func() {
		defer panicked.Recover()
		ret, err := fn()
		finished <- &channelValue{ret, err}
	}()
	select {
	case result := <-finished:
		if result.err != nil {
			limiter.onError(result.err)
		} else {
			limiter.onSuccess()
		}
		return result.ret, result.err
	case <-timeout.Done():
		limiter.onTimeout()
		return nil, &TimeoutError{name: limiter.name}
	case v := <-panicked.Caught():
		limiter.onPanic(v)
		panic(v)
	}
}

func (limiter *timeLimiter) onSuccess() {
	limiter.metrics.successIncrement()
	limiter.eventListener.consumeEvent(newSuccessEvent(limiter.name))
}

func (limiter *timeLimiter) onError(err error) {
	limiter.metrics.errorIncrement()
	limiter.eventListener.consumeEvent(newErrorEvent(limiter.name, err))
}

func (limiter *timeLimiter) onTimeout() {
	limiter.metrics.timeoutIncrement()
	limiter.eventListener.consumeEvent(newTimeoutEvent(limiter.name))
}

func (limiter *timeLimiter) onPanic(v any) {
	limiter.metrics.panicIncrement()
	limiter.eventListener.consumeEvent(newPanicEvent(limiter.name, v))
}

type channelValue struct {
	ret any
	err error
}

type TimeoutError struct {
	name string
}

func (e *TimeoutError) Error() string {
	return fmt.Sprintf("TimeLimiter '%s' recorded a timeout exception.", e.name)
}
