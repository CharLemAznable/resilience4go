package circuitbreaker

import (
	"fmt"
	"github.com/CharLemAznable/resilience4go/utils"
	"time"
)

type SlidingWindowType string

const (
	TimeBased  SlidingWindowType = "TIME_BASED"
	CountBased SlidingWindowType = "COUNT_BASED"
)

type Config struct {
	slidingWindowType                            SlidingWindowType
	slidingWindowSize                            int64
	minimumNumberOfCalls                         int64
	failureRateThreshold                         float64
	slowCallRateThreshold                        float64
	slowCallDurationThreshold                    time.Duration
	failureResultPredicate                       func(any, error) bool
	automaticTransitionFromOpenToHalfOpenEnabled bool
	waitIntervalFunctionInOpenState              func(int64) time.Duration
	permittedNumberOfCallsInHalfOpenState        int64
	maxWaitDurationInHalfOpenState               time.Duration
}

func (config *Config) String() string {
	return fmt.Sprintf(
		"CircuitBreakerConfig"+
			" {slidingWindowType=%s, slidingWindowSize=%d, minimumNumberOfCalls=%d"+
			", failureRateThreshold=%f, slowCallRateThreshold=%f, slowCallDurationThreshold=%v"+
			", failureResultPredicate %T[%v]"+
			", automaticTransitionFromOpenToHalfOpenEnabled=%t"+
			", waitIntervalFunctionInOpenState %T[%v]"+
			", permittedNumberOfCallsInHalfOpenState=%d, maxWaitDurationInHalfOpenState=%v}",
		config.slidingWindowType, config.slidingWindowSize, config.minimumNumberOfCalls,
		config.failureRateThreshold, config.slowCallRateThreshold, config.slowCallDurationThreshold,
		config.failureResultPredicate, any(config.failureResultPredicate),
		config.automaticTransitionFromOpenToHalfOpenEnabled,
		config.waitIntervalFunctionInOpenState, any(config.waitIntervalFunctionInOpenState),
		config.permittedNumberOfCallsInHalfOpenState, config.maxWaitDurationInHalfOpenState)
}

func (config *Config) failureResultPredicateFn(ret any, err error) bool {
	if config.failureResultPredicate != nil {
		return config.failureResultPredicate(ret, err)
	}
	return DefaultFailureResultPredicate(ret, err)
}

func (config *Config) waitIntervalFunctionInOpenStateFn(attempts int64) time.Duration {
	if config.waitIntervalFunctionInOpenState != nil {
		return config.waitIntervalFunctionInOpenState(attempts)
	}
	return DefaultWaitIntervalFunctionInOpenState(attempts)
}

type ConfigBuilder func(*Config)

func WithSlidingWindow(slidingWindowType SlidingWindowType, slidingWindowSize int64, minimumNumberOfCalls int64) ConfigBuilder {
	return func(config *Config) {
		if slidingWindowSize < 1 {
			panic("slidingWindowSize must be greater than 0")
		}
		if minimumNumberOfCalls < 1 {
			panic("minimumNumberOfCalls must be greater than 0")
		}
		config.slidingWindowType = slidingWindowType
		config.slidingWindowSize = slidingWindowSize
		if CountBased == slidingWindowType {
			config.minimumNumberOfCalls = utils.Min(minimumNumberOfCalls, slidingWindowSize)
		} else {
			config.minimumNumberOfCalls = minimumNumberOfCalls
		}
	}
}

func WithFailureRateThreshold(failureRateThreshold float64) ConfigBuilder {
	return func(config *Config) {
		config.failureRateThreshold = failureRateThreshold
	}
}

func WithSlowCallRateThreshold(slowCallRateThreshold float64) ConfigBuilder {
	return func(config *Config) {
		config.slowCallRateThreshold = slowCallRateThreshold
	}
}

func WithSlowCallDurationThreshold(slowCallDurationThreshold time.Duration) ConfigBuilder {
	return func(config *Config) {
		config.slowCallDurationThreshold = slowCallDurationThreshold
	}
}

func WithFailureResultPredicate(predicate func(any, error) bool) ConfigBuilder {
	return func(config *Config) {
		config.failureResultPredicate = predicate
	}
}

func WithAutomaticTransitionFromOpenToHalfOpenEnabled(enabled bool) ConfigBuilder {
	return func(config *Config) {
		config.automaticTransitionFromOpenToHalfOpenEnabled = enabled
	}
}

func WithWaitIntervalFunctionInOpenState(function func(int64) time.Duration) ConfigBuilder {
	return func(config *Config) {
		config.waitIntervalFunctionInOpenState = function
	}
}

func WithPermittedNumberOfCallsInHalfOpenState(permittedNumberOfCallsInHalfOpenState int64) ConfigBuilder {
	return func(config *Config) {
		config.permittedNumberOfCallsInHalfOpenState = permittedNumberOfCallsInHalfOpenState
	}
}

func WithMaxWaitDurationInHalfOpenState(maxWaitDurationInHalfOpenState time.Duration) ConfigBuilder {
	return func(config *Config) {
		config.maxWaitDurationInHalfOpenState = maxWaitDurationInHalfOpenState
	}
}

const DefaultSlidingWindowType = CountBased
const DefaultSlidingWindowSize int64 = 100
const DefaultMinimumNumberOfCalls int64 = 100
const DefaultFailureRateThreshold float64 = 50
const DefaultSlowCallRateThreshold float64 = 100
const DefaultSlowCallDurationThreshold = time.Second * 60

func DefaultFailureResultPredicate(_ any, err error) bool {
	return err != nil
}

const DefaultAutomaticTransitionFromOpenToHalfOpenEnabled bool = false
const DefaultWaitDurationInOpenState = time.Second * 60

func DefaultWaitIntervalFunctionInOpenState(_ int64) time.Duration {
	return DefaultWaitDurationInOpenState
}

const DefaultPermittedNumberOfCallsInHalfOpenState int64 = 10
const DefaultMaxWaitDurationInHalfOpenState time.Duration = 0

func defaultConfig() *Config {
	return &Config{
		slidingWindowType:                            DefaultSlidingWindowType,
		slidingWindowSize:                            DefaultSlidingWindowSize,
		minimumNumberOfCalls:                         DefaultMinimumNumberOfCalls,
		failureRateThreshold:                         DefaultFailureRateThreshold,
		slowCallRateThreshold:                        DefaultSlowCallRateThreshold,
		slowCallDurationThreshold:                    DefaultSlowCallDurationThreshold,
		failureResultPredicate:                       DefaultFailureResultPredicate,
		automaticTransitionFromOpenToHalfOpenEnabled: DefaultAutomaticTransitionFromOpenToHalfOpenEnabled,
		waitIntervalFunctionInOpenState:              DefaultWaitIntervalFunctionInOpenState,
		permittedNumberOfCallsInHalfOpenState:        DefaultPermittedNumberOfCallsInHalfOpenState,
		maxWaitDurationInHalfOpenState:               DefaultMaxWaitDurationInHalfOpenState,
	}
}
