package bulkhead

import (
	"context"
	"fmt"
	"golang.org/x/sync/semaphore"
)

type Bulkhead interface {
	EventProcessor() EventProcessor

	acquire() error
	release()
}

func NewBulkhead(name string, configs ...ConfigBuilder) Bulkhead {
	config := defaultConfig()
	for _, cfg := range configs {
		cfg(config)
	}
	return &semaphoreBulkhead{
		name:           name,
		config:         config,
		rootContext:    context.Background(),
		semaphore:      semaphore.NewWeighted(config.maxConcurrentCalls),
		eventProcessor: newEventProcessor(),
	}
}

type semaphoreBulkhead struct {
	name           string
	config         *Config
	rootContext    context.Context
	semaphore      *semaphore.Weighted
	eventProcessor EventProcessor
}

func (bulkhead *semaphoreBulkhead) EventProcessor() EventProcessor {
	return bulkhead.eventProcessor
}

func (bulkhead *semaphoreBulkhead) acquire() error {
	permitted := func() bool {
		timeout, cancelFn := context.WithTimeout(
			bulkhead.rootContext,
			bulkhead.config.maxWaitDuration)
		defer cancelFn()
		return bulkhead.semaphore.Acquire(timeout, 1) == nil
	}()

	bulkhead.eventProcessor.consumeEvent(func() Event {
		if permitted {
			return newPermittedEvent(bulkhead.name)
		} else {
			return newRejectedEvent(bulkhead.name)
		}
	}())

	if permitted {
		return nil
	}
	return &bulkheadError{name: bulkhead.name}
}

func (bulkhead *semaphoreBulkhead) release() {
	bulkhead.semaphore.Release(1)
	bulkhead.eventProcessor.consumeEvent(newFinishedEvent(bulkhead.name))
}

type bulkheadError struct {
	name string
}

func (e *bulkheadError) Error() string {
	return fmt.Sprintf("Bulkhead '%s' is full and does not permit further calls", e.name)
}
