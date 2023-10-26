package bulkhead_test

import (
	"fmt"
	"github.com/CharLemAznable/resilience4go/bulkhead"
	"testing"
	"time"
)

func TestConfig_String(t *testing.T) {
	config := &bulkhead.Config{}
	bulkhead.WithMaxConcurrentCalls(10)(config)
	bulkhead.WithMaxWaitDuration(time.Second * 5)(config)
	expected := "BulkheadConfig{maxConcurrentCalls=10, maxWaitDuration=5s}"
	result := fmt.Sprintf("%v", config)
	if result != expected {
		t.Errorf("Expected %s, but got %s", expected, result)
	}
}
