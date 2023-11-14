package retry

import (
	"fmt"
	"time"
)

type Config struct {
	maxAttempts            int
	failAfterMaxAttempts   bool
	failureResultPredicate func(any, error) bool
	waitIntervalFunction   func(int) time.Duration
}

func (config *Config) String() string {
	return fmt.Sprintf(
		"RetryConfig"+
			" {maxAttempts=%d, failAfterMaxAttempts=%t"+
			", failureResultPredicate %T[%v], waitIntervalFunction %T[%v]}",
		config.maxAttempts, config.failAfterMaxAttempts,
		config.failureResultPredicate, any(config.failureResultPredicate),
		config.waitIntervalFunction, any(config.waitIntervalFunction))
}

func (config *Config) failureResultPredicateFn(ret any, err error) bool {
	if config.failureResultPredicate != nil {
		return config.failureResultPredicate(ret, err)
	}
	return DefaultFailureResultPredicate(ret, err)
}

func (config *Config) waitIntervalFunctionFn(attempts int) time.Duration {
	if config.waitIntervalFunction != nil {
		return config.waitIntervalFunction(attempts)
	}
	return DefaultWaitIntervalFunction(attempts)
}

type ConfigBuilder func(*Config)

func WithMaxAttempts(maxAttempts int) ConfigBuilder {
	return func(config *Config) {
		config.maxAttempts = maxAttempts
	}
}

func WithFailAfterMaxAttempts(failAfterMaxAttempts bool) ConfigBuilder {
	return func(config *Config) {
		config.failAfterMaxAttempts = failAfterMaxAttempts
	}
}

func WithFailureResultPredicate(predicate func(any, error) bool) ConfigBuilder {
	return func(config *Config) {
		config.failureResultPredicate = predicate
	}
}

func WithWaitIntervalFunction(function func(int) time.Duration) ConfigBuilder {
	return func(config *Config) {
		config.waitIntervalFunction = function
	}
}

const DefaultMaxAttempts int = 3
const DefaultFailAfterMaxAttempts bool = false

func DefaultFailureResultPredicate(_ any, err error) bool {
	return err != nil
}

const DefaultWaitDuration = time.Millisecond * 500

func DefaultWaitIntervalFunction(_ int) time.Duration {
	return DefaultWaitDuration
}

func defaultConfig() *Config {
	return &Config{
		maxAttempts:            DefaultMaxAttempts,
		failAfterMaxAttempts:   DefaultFailAfterMaxAttempts,
		failureResultPredicate: DefaultFailureResultPredicate,
		waitIntervalFunction:   DefaultWaitIntervalFunction,
	}
}
