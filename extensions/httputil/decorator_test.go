package httputil_test

import (
	"bufio"
	"errors"
	"github.com/CharLemAznable/resilience4go/bulkhead"
	"github.com/CharLemAznable/resilience4go/circuitbreaker"
	"github.com/CharLemAznable/resilience4go/extensions/httputil"
	"github.com/CharLemAznable/resilience4go/ratelimiter"
	"github.com/CharLemAznable/resilience4go/retry"
	"github.com/CharLemAznable/resilience4go/timelimiter"
	"io"
	"net/http"
	"strings"
	"testing"
)

func TestDecorateRoundTripper(t *testing.T) {
	decoratedRoundTripper := httputil.
		OfRoundTripper(httputil.RoundTripperFunc(
			func(req *http.Request) (*http.Response, error) {
				return nil, errors.New("error")
			})).
		WithBulkhead(bulkhead.NewBulkhead("test")).
		WithTimeLimiter(timelimiter.NewTimeLimiter("test")).
		WithRateLimiter(ratelimiter.NewRateLimiter("test")).
		WithCircuitBreaker(circuitbreaker.NewCircuitBreaker("test")).
		WithRetry(retry.NewRetry("test")).
		WithFallback(func(err error) (*http.Response, error) {
			return http.ReadResponse(bufio.NewReader(strings.NewReader(
				"HTTP/1.1 200 OK\r\n"+
					"\r\n"+
					"fallback")), &http.Request{Method: "GET"})
		}).
		Decorate()

	if decoratedRoundTripper == nil {
		t.Error("Expected non-nil decoratedRoundTripper")
	}
	resp, err := decoratedRoundTripper.RoundTrip(nil)
	if resp.StatusCode != 200 {
		t.Errorf("Expected resp 200, but got '%d'", resp.StatusCode)
	}
	if err != nil {
		t.Errorf("Expected error is nil, but got '%v'", err)
	}
	defer func() { _ = resp.Body.Close() }()
	body, err := io.ReadAll(resp.Body)
	if string(body) != "fallback" {
		t.Errorf("Expected ret is 'fallback', but got '%v'", string(body))
	}
	if err != nil {
		t.Errorf("Expected error is nil, but got '%v'", err)
	}
}
