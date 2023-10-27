package timelimiter_test

import (
	"fmt"
	"github.com/CharLemAznable/resilience4go/timelimiter"
	"testing"
	"time"
)

func TestConfig_String(t *testing.T) {
	config := &timelimiter.Config{}
	timelimiter.WithTimeoutDuration(time.Second * 5)(config)
	expected := "TimeLimiterConfig{timeoutDuration=5s}"
	result := fmt.Sprintf("%v", config)
	if result != expected {
		t.Errorf("Expected %s, but got %s", expected, result)
	}
}
