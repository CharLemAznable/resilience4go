package retry

import (
	"fmt"
	"time"
)

type Config struct {
	maxAttempts           int
	failAfterMaxAttempts  bool
	recordResultPredicate func(any, error) bool
	waitIntervalFunction  func(int) time.Duration
}

func (config *Config) String() string {
	return fmt.Sprintf(
		"RetryConfig"+
			" {maxAttempts=%d, failAfterMaxAttempts=%t"+
			", recordResultPredicate %T[%v], waitIntervalFunction %T[%v]}",
		config.maxAttempts, config.failAfterMaxAttempts,
		config.recordResultPredicate, any(config.recordResultPredicate),
		config.waitIntervalFunction, any(config.waitIntervalFunction))
}

func (config *Config) recordResultPredicateFn(ret any, err error) bool {
	if config.recordResultPredicate != nil {
		return config.recordResultPredicate(ret, err)
	}
	return DefaultRecordResultPredicate(ret, err)
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

func WithRecordResultPredicate(predicate func(any, error) bool) ConfigBuilder {
	return func(config *Config) {
		config.recordResultPredicate = predicate
	}
}

func WithWaitIntervalFunction(function func(int) time.Duration) ConfigBuilder {
	return func(config *Config) {
		config.waitIntervalFunction = function
	}
}

const DefaultMaxAttempts int = 3
const DefaultFailAfterMaxAttempts bool = false

func DefaultRecordResultPredicate(_ any, err error) bool {
	return err != nil
}

const DefaultWaitDuration = time.Millisecond * 500

func DefaultWaitIntervalFunction(_ int) time.Duration {
	return DefaultWaitDuration
}

func defaultConfig() *Config {
	return &Config{
		maxAttempts:           DefaultMaxAttempts,
		failAfterMaxAttempts:  DefaultFailAfterMaxAttempts,
		recordResultPredicate: DefaultRecordResultPredicate,
		waitIntervalFunction:  DefaultWaitIntervalFunction,
	}
}
