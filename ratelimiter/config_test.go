package ratelimiter_test

import (
	"fmt"
	"github.com/CharLemAznable/resilience4go/ratelimiter"
	"testing"
	"time"
)

func TestConfig_String(t *testing.T) {
	config := &ratelimiter.Config{}
	ratelimiter.WithTimeoutDuration(time.Second * 10)(config)
	ratelimiter.WithLimitRefreshPeriod(time.Nanosecond * 100)(config)
	ratelimiter.WithLimitForPeriod(10)(config)
	expected := "RateLimiterConfig{timeoutDuration=10s, limitRefreshPeriod=100ns, limitForPeriod=10}"
	result := fmt.Sprintf("%v", config)
	if result != expected {
		t.Errorf("Expected %s, but got %s", expected, result)
	}
}
