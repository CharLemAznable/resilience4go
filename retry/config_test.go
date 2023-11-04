package retry_test

import (
	"fmt"
	"github.com/CharLemAznable/resilience4go/retry"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestConfig_String(t *testing.T) {
	config := &retry.Config{}
	retry.WithMaxAttempts(5)(config)
	retry.WithFailAfterMaxAttempts(true)(config)
	recordResultPredicate := func(ret any, err error) bool {
		return ret == nil || err != nil
	}
	retry.WithRecordResultPredicate(recordResultPredicate)(config)
	waitIntervalFunction := func(_ int) time.Duration {
		return time.Second * 30
	}
	retry.WithWaitIntervalFunction(waitIntervalFunction)(config)
	expected := fmt.Sprintf("RetryConfig"+
		" {maxAttempts=5, failAfterMaxAttempts=true"+
		", recordResultPredicate %T[%v]"+
		", waitIntervalFunction %T[%v]}",
		recordResultPredicate, any(recordResultPredicate),
		waitIntervalFunction, any(waitIntervalFunction))
	result := fmt.Sprintf("%v", config)
	assert.Equal(t, expected, result)
}
