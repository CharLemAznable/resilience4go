package retry_test

import (
	"fmt"
	"github.com/CharLemAznable/resilience4go/retry"
	"testing"
	"time"
)

func TestConfig_String(t *testing.T) {
	config := &retry.Config{}
	retry.WithMaxAttempts(5)(config)
	retry.WithFailAfterMaxAttempts(true)(config)
	failureResultPredicate := func(ret any, err error) bool {
		return ret == nil || err != nil
	}
	retry.WithFailureResultPredicate(failureResultPredicate)(config)
	waitIntervalFunction := func(_ int) time.Duration {
		return time.Second * 30
	}
	retry.WithWaitIntervalFunction(waitIntervalFunction)(config)
	expected := fmt.Sprintf("RetryConfig"+
		" {maxAttempts=5, failAfterMaxAttempts=true"+
		", failureResultPredicate %T[%v]"+
		", waitIntervalFunction %T[%v]}",
		failureResultPredicate, any(failureResultPredicate),
		waitIntervalFunction, any(waitIntervalFunction))
	result := fmt.Sprintf("%v", config)
	if result != expected {
		t.Errorf("Expected config string '%s', but got '%s'", expected, result)
	}
}
