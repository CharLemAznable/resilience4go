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
	return OfConsumer(bulkhead.DecorateConsumer(entry, consumer.Consumer))
}

func (consumer *DecorateConsumer[T]) WhenFull(fn func(T) error) *DecorateConsumer[T] {
	return OfConsumer(fallback.DecorateConsumerByType[T, *bulkhead.FullError](consumer.Consumer, fn))
}

func (consumer *DecorateConsumer[T]) WithTimeLimiter(entry timelimiter.TimeLimiter) *DecorateConsumer[T] {
	return OfConsumer(timelimiter.DecorateConsumer(entry, consumer.Consumer))
}

func (consumer *DecorateConsumer[T]) WhenTimeout(fn func(T) error) *DecorateConsumer[T] {
	return OfConsumer(fallback.DecorateConsumerByType[T, *timelimiter.TimeoutError](consumer.Consumer, fn))
}

func (consumer *DecorateConsumer[T]) WithRateLimiter(entry ratelimiter.RateLimiter) *DecorateConsumer[T] {
	return OfConsumer(ratelimiter.DecorateConsumer(entry, consumer.Consumer))
}

func (consumer *DecorateConsumer[T]) WhenOverRate(fn func(T) error) *DecorateConsumer[T] {
	return OfConsumer(fallback.DecorateConsumerByType[T, *ratelimiter.NotPermittedError](consumer.Consumer, fn))
}

func (consumer *DecorateConsumer[T]) WithCircuitBreaker(entry circuitbreaker.CircuitBreaker) *DecorateConsumer[T] {
	return OfConsumer(circuitbreaker.DecorateConsumer(entry, consumer.Consumer))
}

func (consumer *DecorateConsumer[T]) WhenOverLoad(fn func(T) error) *DecorateConsumer[T] {
	return OfConsumer(fallback.DecorateConsumerByType[T, *circuitbreaker.NotPermittedError](consumer.Consumer, fn))
}

func (consumer *DecorateConsumer[T]) WithRetry(entry retry.Retry) *DecorateConsumer[T] {
	return OfConsumer(retry.DecorateConsumer(entry, consumer.Consumer))
}

func (consumer *DecorateConsumer[T]) WhenMaxRetries(fn func(T) error) *DecorateConsumer[T] {
	return OfConsumer(fallback.DecorateConsumerByType[T, *retry.MaxRetriesExceeded](consumer.Consumer, fn))
}

func (consumer *DecorateConsumer[T]) WithFallback(
	fn func(T) error, predicate func(T, error, any) bool) *DecorateConsumer[T] {
	return OfConsumer(fallback.DecorateConsumer(consumer.Consumer,
		func(ctx fallback.Context[T, any, error]) error { return fn(ctx.Param()) },
		func(ctx fallback.Context[T, any, error]) (bool, fallback.Context[T, any, error]) {
			return predicate(ctx.Param(), ctx.Err(), ctx.Panic()), ctx
		}))
}

func (consumer *DecorateConsumer[T]) Decorate() consumer.Consumer[T] {
	return consumer.Consumer
}
