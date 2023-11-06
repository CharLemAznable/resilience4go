package bulkhead

import (
	"context"
	"fmt"
	"golang.org/x/sync/semaphore"
)

type Bulkhead interface {
	Name() string
	Metrics() Metrics
	EventListener() EventListener

	acquire() error
	release()
}

func NewBulkhead(name string, configs ...ConfigBuilder) Bulkhead {
	config := defaultConfig()
	for _, cfg := range configs {
		cfg(config)
	}
	return &semaphoreBulkhead{
		name:          name,
		config:        config,
		rootContext:   context.Background(),
		semaphore:     semaphore.NewWeighted(config.maxConcurrentCalls),
		metrics:       newMetric(config.maxConcurrentCalls),
		eventListener: newEventListener(),
	}
}

type semaphoreBulkhead struct {
	name          string
	config        *Config
	rootContext   context.Context
	semaphore     *semaphore.Weighted
	metrics       Metrics
	eventListener EventListener
}

func (bulkhead *semaphoreBulkhead) Name() string {
	return bulkhead.name
}

func (bulkhead *semaphoreBulkhead) Metrics() Metrics {
	return bulkhead.metrics
}

func (bulkhead *semaphoreBulkhead) EventListener() EventListener {
	return bulkhead.eventListener
}

func (bulkhead *semaphoreBulkhead) acquire() error {
	permitted := func() bool {
		timeout, cancelFn := context.WithTimeout(
			bulkhead.rootContext,
			bulkhead.config.maxWaitDuration)
		defer cancelFn()
		return bulkhead.semaphore.Acquire(timeout, 1) == nil
	}()

	bulkhead.eventListener.consumeEvent(func() Event {
		if permitted {
			return newPermittedEvent(bulkhead.name)
		} else {
			return newRejectedEvent(bulkhead.name)
		}
	}())

	if permitted {
		bulkhead.metrics.acquire(1)
		return nil
	}
	return &FullError{name: bulkhead.name}
}

func (bulkhead *semaphoreBulkhead) release() {
	bulkhead.semaphore.Release(1)
	bulkhead.metrics.release(1)
	bulkhead.eventListener.consumeEvent(newFinishedEvent(bulkhead.name))
}

type FullError struct {
	name string
}

func (e *FullError) Error() string {
	return fmt.Sprintf("Bulkhead '%s' is full and does not permit further calls", e.name)
}
