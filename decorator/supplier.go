package decorator

import (
	"github.com/CharLemAznable/gofn/supplier"
	"github.com/CharLemAznable/resilience4go/bulkhead"
	"github.com/CharLemAznable/resilience4go/circuitbreaker"
	"github.com/CharLemAznable/resilience4go/fallback"
	"github.com/CharLemAznable/resilience4go/ratelimiter"
	"github.com/CharLemAznable/resilience4go/retry"
	"github.com/CharLemAznable/resilience4go/timelimiter"
)

func OfSupplier[T any](supplier supplier.Supplier[T]) *DecorateSupplier[T] {
	return &DecorateSupplier[T]{supplier}
}

type DecorateSupplier[T any] struct {
	supplier.Supplier[T]
}

func (function *DecorateSupplier[T]) WithBulkhead(entry bulkhead.Bulkhead) *DecorateSupplier[T] {
	return OfSupplier(bulkhead.DecorateSupplier(entry, function.Supplier))
}

func (function *DecorateSupplier[T]) WhenFull(fn func(T, *bulkhead.FullError) (T, error)) *DecorateSupplier[T] {
	return OfSupplier(fallback.DecorateSupplierDefault(function.Supplier, fn))
}

func (function *DecorateSupplier[T]) WithTimeLimiter(entry timelimiter.TimeLimiter) *DecorateSupplier[T] {
	return OfSupplier(timelimiter.DecorateSupplier(entry, function.Supplier))
}

func (function *DecorateSupplier[T]) WhenTimeout(fn func(T, *timelimiter.TimeoutError) (T, error)) *DecorateSupplier[T] {
	return OfSupplier(fallback.DecorateSupplierDefault(function.Supplier, fn))
}

func (function *DecorateSupplier[T]) WithRateLimiter(entry ratelimiter.RateLimiter) *DecorateSupplier[T] {
	return OfSupplier(ratelimiter.DecorateSupplier(entry, function.Supplier))
}

func (function *DecorateSupplier[T]) WhenOverRate(fn func(T, *ratelimiter.NotPermittedError) (T, error)) *DecorateSupplier[T] {
	return OfSupplier(fallback.DecorateSupplierDefault(function.Supplier, fn))
}

func (function *DecorateSupplier[T]) WithCircuitBreaker(entry circuitbreaker.CircuitBreaker) *DecorateSupplier[T] {
	return OfSupplier(circuitbreaker.DecorateSupplier(entry, function.Supplier))
}

func (function *DecorateSupplier[T]) WhenOverLoad(fn func(T, *circuitbreaker.NotPermittedError) (T, error)) *DecorateSupplier[T] {
	return OfSupplier(fallback.DecorateSupplierDefault(function.Supplier, fn))
}

func (function *DecorateSupplier[T]) WithRetry(entry retry.Retry) *DecorateSupplier[T] {
	return OfSupplier(retry.DecorateSupplier(entry, function.Supplier))
}

func (function *DecorateSupplier[T]) WhenMaxRetries(fn func(T, *retry.MaxRetriesExceeded) (T, error)) *DecorateSupplier[T] {
	return OfSupplier(fallback.DecorateSupplierDefault(function.Supplier, fn))
}

func (function *DecorateSupplier[T]) WithFallback(fn func(T, error) (T, error)) *DecorateSupplier[T] {
	return OfSupplier(fallback.DecorateSupplierDefault(function.Supplier, fn))
}

func (function *DecorateSupplier[T]) Decorate() supplier.Supplier[T] {
	return function.Supplier
}
