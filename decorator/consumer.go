package decorator

import (
	"github.com/CharLemAznable/gofn/consumer"
	"github.com/CharLemAznable/resilience4go/bulkhead"
	"github.com/CharLemAznable/resilience4go/circuitbreaker"
	"github.com/CharLemAznable/resilience4go/fallback"
	"github.com/CharLemAznable/resilience4go/ratelimiter"
	"github.com/CharLemAznable/resilience4go/retry"
	"github.com/CharLemAznable/resilience4go/timelimiter"
)

func OfConsumer[T any](consumer consumer.Consumer[T]) *DecorateConsumer[T] {
	return &DecorateConsumer[T]{consumer}
}

type DecorateConsumer[T any] struct {
	consumer.Consumer[T]
}

func (consumer *DecorateConsumer[T]) WithBulkhead(entry bulkhead.Bulkhead) *DecorateConsumer[T] {
	return &DecorateConsumer[T]{bulkhead.DecorateConsumer(entry, consumer.Consumer)}
}

func (consumer *DecorateConsumer[T]) WithTimeLimiter(entry timelimiter.TimeLimiter) *DecorateConsumer[T] {
	return &DecorateConsumer[T]{timelimiter.DecorateConsumer(entry, consumer.Consumer)}
}

func (consumer *DecorateConsumer[T]) WithRateLimiter(entry ratelimiter.RateLimiter) *DecorateConsumer[T] {
	return &DecorateConsumer[T]{ratelimiter.DecorateConsumer(entry, consumer.Consumer)}
}

func (consumer *DecorateConsumer[T]) WithCircuitBreaker(entry circuitbreaker.CircuitBreaker) *DecorateConsumer[T] {
	return &DecorateConsumer[T]{circuitbreaker.DecorateConsumer(entry, consumer.Consumer)}
}

func (consumer *DecorateConsumer[T]) WithRetry(entry retry.Retry) *DecorateConsumer[T] {
	return &DecorateConsumer[T]{retry.DecorateConsumer(entry, consumer.Consumer)}
}

func (consumer *DecorateConsumer[T]) WithFallback(fn func(error) error) *DecorateConsumer[T] {
	return &DecorateConsumer[T]{fallback.DecorateConsumer(consumer.Consumer, fn)}
}

func (consumer *DecorateConsumer[T]) Decorate() consumer.Consumer[T] {
	return consumer.Consumer
}
