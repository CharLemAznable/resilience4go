package ratelimiter

import (
	"fmt"
	"time"
)

type Config struct {
	timeoutDuration    time.Duration
	limitRefreshPeriod time.Duration
	limitForPeriod     int64
}

func (config *Config) String() string {
	return fmt.Sprintf(
		"RateLimiterConfig{timeoutDuration=%v, limitRefreshPeriod=%v, limitForPeriod=%d}",
		config.timeoutDuration, config.limitRefreshPeriod, config.limitForPeriod)
}

type ConfigBuilder func(*Config)

func WithTimeoutDuration(timeoutDuration time.Duration) ConfigBuilder {
	return func(config *Config) {
		config.timeoutDuration = timeoutDuration
	}
}

func WithLimitRefreshPeriod(limitRefreshPeriod time.Duration) ConfigBuilder {
	return func(config *Config) {
		config.limitRefreshPeriod = limitRefreshPeriod
	}
}

func WithLimitForPeriod(limitForPeriod int64) ConfigBuilder {
	return func(config *Config) {
		config.limitForPeriod = limitForPeriod
	}
}

const DefaultTimeoutDuration = time.Second * 5
const DefaultLimitRefreshPeriod = time.Nanosecond * 500
const DefaultLimitForPeriod int64 = 50

func defaultConfig() *Config {
	return &Config{
		timeoutDuration:    DefaultTimeoutDuration,
		limitRefreshPeriod: DefaultLimitRefreshPeriod,
		limitForPeriod:     DefaultLimitForPeriod,
	}
}
