package ratelimiter_test

import (
	"github.com/CharLemAznable/resilience4go/ratelimiter"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestDecorateRunnable(t *testing.T) {
	// 创建一个可运行的函数
	fn := func() error {
		time.Sleep(time.Millisecond * 500)
		return nil
	}

	// 创建一个Bulkhead的mock对象
	rl := ratelimiter.NewRateLimiter("test",
		ratelimiter.WithTimeoutDuration(time.Second),
		ratelimiter.WithLimitRefreshPeriod(time.Second*2),
		ratelimiter.WithLimitForPeriod(2))

	// 调用DecorateRunnable函数
	decoratedFn := ratelimiter.DecorateRunnable(rl, fn)

	errChan := make(chan error, 1)

	for i := 0; i < 3; i++ {
		go func() {
			err := decoratedFn()
			if err != nil {
				errChan <- err
			}
		}()
	}

	fullErr, ok := (<-errChan).(*ratelimiter.NotPermittedError)
	assert.True(t, ok)
	assert.Equal(t, "RateLimiter 'test' does not permit further calls", fullErr.Error())
}

func TestDecorateSupplier(t *testing.T) {
	// 创建一个可运行的函数
	fn := func() (string, error) {
		time.Sleep(time.Millisecond * 500)
		return "error", nil
	}

	// 创建一个Bulkhead的mock对象
	rl := ratelimiter.NewRateLimiter("test",
		ratelimiter.WithTimeoutDuration(time.Second),
		ratelimiter.WithLimitRefreshPeriod(time.Second*2),
		ratelimiter.WithLimitForPeriod(2))

	// 调用DecorateRunnable函数
	decoratedFn := ratelimiter.DecorateSupplier(rl, fn)

	errChan := make(chan error, 1)

	for i := 0; i < 3; i++ {
		go func() {
			_, err := decoratedFn()
			if err != nil {
				errChan <- err
			}
		}()
	}

	fullErr, ok := (<-errChan).(*ratelimiter.NotPermittedError)
	assert.True(t, ok)
	assert.Equal(t, "RateLimiter 'test' does not permit further calls", fullErr.Error())
}

func TestDecorateConsumer(t *testing.T) {
	// 创建一个可运行的函数
	fn := func(str string) error {
		time.Sleep(time.Millisecond * 500)
		return nil
	}

	// 创建一个Bulkhead的mock对象
	rl := ratelimiter.NewRateLimiter("test",
		ratelimiter.WithTimeoutDuration(time.Second),
		ratelimiter.WithLimitRefreshPeriod(time.Second*2),
		ratelimiter.WithLimitForPeriod(2))

	// 调用DecorateRunnable函数
	decoratedFn := ratelimiter.DecorateConsumer(rl, fn)

	errChan := make(chan error, 1)

	for i := 0; i < 3; i++ {
		go func() {
			err := decoratedFn("error")
			if err != nil {
				errChan <- err
			}
		}()
	}

	fullErr, ok := (<-errChan).(*ratelimiter.NotPermittedError)
	assert.True(t, ok)
	assert.Equal(t, "RateLimiter 'test' does not permit further calls", fullErr.Error())
}

func TestDecorateFunction(t *testing.T) {
	// 创建一个可运行的函数
	fn := func(str string) (string, error) {
		time.Sleep(time.Millisecond * 500)
		return str, nil
	}

	// 创建一个Bulkhead的mock对象
	rl := ratelimiter.NewRateLimiter("test",
		ratelimiter.WithTimeoutDuration(time.Second),
		ratelimiter.WithLimitRefreshPeriod(time.Second*2),
		ratelimiter.WithLimitForPeriod(2))

	// 调用DecorateRunnable函数
	decoratedFn := ratelimiter.DecorateFunction(rl, fn)

	errChan := make(chan error, 1)

	for i := 0; i < 3; i++ {
		go func() {
			_, err := decoratedFn("error")
			if err != nil {
				errChan <- err
			}
		}()
	}

	fullErr, ok := (<-errChan).(*ratelimiter.NotPermittedError)
	assert.True(t, ok)
	assert.Equal(t, "RateLimiter 'test' does not permit further calls", fullErr.Error())
}
