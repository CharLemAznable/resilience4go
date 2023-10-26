package bulkhead

import (
	"fmt"
	"time"
)

type Config struct {
	maxConcurrentCalls int64
	maxWaitDuration    time.Duration
}

func (config *Config) String() string {
	return fmt.Sprintf(
		"BulkheadConfig{maxConcurrentCalls=%d, maxWaitDuration=%v}",
		config.maxConcurrentCalls, config.maxWaitDuration)
}

type ConfigBuilder func(*Config)

func WithMaxConcurrentCalls(maxConcurrentCalls int64) ConfigBuilder {
	return func(config *Config) {
		config.maxConcurrentCalls = maxConcurrentCalls
	}
}

func WithMaxWaitDuration(maxWaitDuration time.Duration) ConfigBuilder {
	return func(config *Config) {
		config.maxWaitDuration = maxWaitDuration
	}
}

const DefaultMaxConcurrentCalls int64 = 25
const DefaultMaxWaitDuration time.Duration = 0

func defaultConfig() *Config {
	return &Config{
		maxConcurrentCalls: DefaultMaxConcurrentCalls,
		maxWaitDuration:    DefaultMaxWaitDuration,
	}
}
