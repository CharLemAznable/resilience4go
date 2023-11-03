package circuitbreaker_test

import (
	"fmt"
	"github.com/CharLemAznable/resilience4go/circuitbreaker"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestConfig_String(t *testing.T) {
	config := &circuitbreaker.Config{}
	circuitbreaker.WithSlidingWindow(circuitbreaker.TimeBased, 50, 50)(config)
	circuitbreaker.WithFailureRateThreshold(75)(config)
	circuitbreaker.WithSlowCallRateThreshold(90)(config)
	circuitbreaker.WithSlowCallDurationThreshold(time.Second * 30)(config)
	recordResultPredicate := func(ret any, err error) bool {
		return ret == nil || err != nil
	}
	circuitbreaker.WithRecordResultPredicate(recordResultPredicate)(config)
	circuitbreaker.WithAutomaticTransitionFromOpenToHalfOpenEnabled(true)(config)
	waitIntervalFunctionInOpenState := func(_ int64) time.Duration {
		return time.Second * 30
	}
	circuitbreaker.WithWaitIntervalFunctionInOpenState(waitIntervalFunctionInOpenState)(config)
	circuitbreaker.WithPermittedNumberOfCallsInHalfOpenState(5)(config)
	circuitbreaker.WithMaxWaitDurationInHalfOpenState(time.Second)(config)
	expected := fmt.Sprintf("CircuitBreakerConfig"+
		" {slidingWindowType=TIME_BASED, slidingWindowSize=50, minimumNumberOfCalls=50"+
		", failureRateThreshold=75.000000, slowCallRateThreshold=90.000000, slowCallDurationThreshold=30s"+
		", recordResultPredicate %T[%v]"+
		", automaticTransitionFromOpenToHalfOpenEnabled=true"+
		", waitIntervalFunctionInOpenState %T[%v]"+
		", permittedNumberOfCallsInHalfOpenState=5, maxWaitDurationInHalfOpenState=1s}",
		recordResultPredicate, any(recordResultPredicate),
		waitIntervalFunctionInOpenState, any(waitIntervalFunctionInOpenState))
	result := fmt.Sprintf("%v", config)
	assert.Equal(t, expected, result)

	assert.PanicsWithValue(t, "slidingWindowSize must be greater than 0", func() {
		circuitbreaker.WithSlidingWindow(circuitbreaker.TimeBased, 0, 50)(config)
	})
	assert.PanicsWithValue(t, "minimumNumberOfCalls must be greater than 0", func() {
		circuitbreaker.WithSlidingWindow(circuitbreaker.TimeBased, 50, 0)(config)
	})

	circuitbreaker.WithSlidingWindow(circuitbreaker.CountBased, 10, 50)(config)
	expected = fmt.Sprintf("CircuitBreakerConfig"+
		" {slidingWindowType=COUNT_BASED, slidingWindowSize=10, minimumNumberOfCalls=10"+
		", failureRateThreshold=75.000000, slowCallRateThreshold=90.000000, slowCallDurationThreshold=30s"+
		", recordResultPredicate %T[%v]"+
		", automaticTransitionFromOpenToHalfOpenEnabled=true"+
		", waitIntervalFunctionInOpenState %T[%v]"+
		", permittedNumberOfCallsInHalfOpenState=5, maxWaitDurationInHalfOpenState=1s}",
		recordResultPredicate, any(recordResultPredicate),
		waitIntervalFunctionInOpenState, any(waitIntervalFunctionInOpenState))
	result = fmt.Sprintf("%v", config)
	assert.Equal(t, expected, result)
}
