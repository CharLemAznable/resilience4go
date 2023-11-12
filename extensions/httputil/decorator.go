package httputil

import (
	"github.com/CharLemAznable/resilience4go/bulkhead"
	"github.com/CharLemAznable/resilience4go/circuitbreaker"
	"github.com/CharLemAznable/resilience4go/fallback"
	"github.com/CharLemAznable/resilience4go/ratelimiter"
	"github.com/CharLemAznable/resilience4go/retry"
	"github.com/CharLemAznable/resilience4go/timelimiter"
	"net/http"
)

type RoundTripperFunc func(*http.Request) (*http.Response, error)

func (fn RoundTripperFunc) RoundTrip(r *http.Request) (*http.Response, error) {
	return fn(r)
}

func OfRoundTripper(rt http.RoundTripper) *DecorateRoundTripper {
	return &DecorateRoundTripper{rt}
}

func OfRoundTripperFunc(fn func(*http.Request) (*http.Response, error)) *DecorateRoundTripper {
	return &DecorateRoundTripper{RoundTripperFunc(fn)}
}

type DecorateRoundTripper struct {
	http.RoundTripper
}

func (rt *DecorateRoundTripper) WithBulkhead(entry bulkhead.Bulkhead) *DecorateRoundTripper {
	return OfRoundTripperFunc(bulkhead.DecorateFunction(entry, rt.RoundTrip))
}

func (rt *DecorateRoundTripper) WithTimeLimiter(entry timelimiter.TimeLimiter) *DecorateRoundTripper {
	return OfRoundTripperFunc(timelimiter.DecorateFunction(entry, rt.RoundTrip))
}

func (rt *DecorateRoundTripper) WithRateLimiter(entry ratelimiter.RateLimiter) *DecorateRoundTripper {
	return OfRoundTripperFunc(ratelimiter.DecorateFunction(entry, rt.RoundTrip))
}

func (rt *DecorateRoundTripper) WithCircuitBreaker(entry circuitbreaker.CircuitBreaker) *DecorateRoundTripper {
	return OfRoundTripperFunc(circuitbreaker.DecorateFunction(entry, rt.RoundTrip))
}

func (rt *DecorateRoundTripper) WithRetry(entry retry.Retry) *DecorateRoundTripper {
	return OfRoundTripperFunc(retry.DecorateFunction(entry, rt.RoundTrip))
}

func (rt *DecorateRoundTripper) WithFallback(fn func(error) (*http.Response, error)) *DecorateRoundTripper {
	return OfRoundTripperFunc(fallback.DecorateFunction(rt.RoundTrip, fn))
}

func (rt *DecorateRoundTripper) Decorate() http.RoundTripper {
	return rt.RoundTripper
}
