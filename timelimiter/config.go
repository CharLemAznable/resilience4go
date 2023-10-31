package timelimiter

import (
	"fmt"
	"time"
)

type Config struct {
	timeoutDuration time.Duration
}

func (config *Config) String() string {
	return fmt.Sprintf(
		"TimeLimiterConfig{timeoutDuration=%v}",
		config.timeoutDuration)
}

type ConfigBuilder func(*Config)

func WithTimeoutDuration(timeoutDuration time.Duration) ConfigBuilder {
	return func(config *Config) {
		config.timeoutDuration = timeoutDuration
	}
}

const DefaultTimeoutDuration = time.Second

func defaultConfig() *Config {
	return &Config{
		timeoutDuration: DefaultTimeoutDuration,
	}
}
