package ratelimiter_test

import (
	"github.com/CharLemAznable/resilience4go/ratelimiter"
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
	decoratedFn := ratelimiter.DecorateCheckedRun(rl, fn)

	errChan := make(chan error, 1)

	for i := 0; i < 3; i++ {
		go func() {
			err := decoratedFn()
			if err != nil {
				errChan <- err
			}
		}()
	}

	err := <-errChan
	fullErr, ok := err.(*ratelimiter.NotPermittedError)
	if !ok {
		t.Errorf("Expected error type *ratelimiter.NotPermittedError, but got '%T'", err)
	} else {
		if fullErr.Error() != "RateLimiter 'test' does not permit further calls" {
			t.Errorf("Expected error message 'RateLimiter 'test' does not permit further calls', but got '%s'", fullErr.Error())
		}
	}
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
	decoratedFn := ratelimiter.DecorateCheckedGet(rl, fn)

	errChan := make(chan error, 1)

	for i := 0; i < 3; i++ {
		go func() {
			_, err := decoratedFn()
			if err != nil {
				errChan <- err
			}
		}()
	}

	err := <-errChan
	fullErr, ok := err.(*ratelimiter.NotPermittedError)
	if !ok {
		t.Errorf("Expected error type *ratelimiter.NotPermittedError, but got '%T'", err)
	} else {
		if fullErr.Error() != "RateLimiter 'test' does not permit further calls" {
			t.Errorf("Expected error message 'RateLimiter 'test' does not permit further calls', but got '%s'", fullErr.Error())
		}
	}
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
	decoratedFn := ratelimiter.DecorateCheckedAccept(rl, fn)

	errChan := make(chan error, 1)

	for i := 0; i < 3; i++ {
		go func() {
			err := decoratedFn("error")
			if err != nil {
				errChan <- err
			}
		}()
	}

	err := <-errChan
	fullErr, ok := err.(*ratelimiter.NotPermittedError)
	if !ok {
		t.Errorf("Expected error type *ratelimiter.NotPermittedError, but got '%T'", err)
	} else {
		if fullErr.Error() != "RateLimiter 'test' does not permit further calls" {
			t.Errorf("Expected error message 'RateLimiter 'test' does not permit further calls', but got '%s'", fullErr.Error())
		}
	}
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
	decoratedFn := ratelimiter.DecorateCheckedApply(rl, fn)

	errChan := make(chan error, 1)

	for i := 0; i < 3; i++ {
		go func() {
			_, err := decoratedFn("error")
			if err != nil {
				errChan <- err
			}
		}()
	}

	err := <-errChan
	fullErr, ok := err.(*ratelimiter.NotPermittedError)
	if !ok {
		t.Errorf("Expected error type *ratelimiter.NotPermittedError, but got '%T'", err)
	} else {
		if fullErr.Error() != "RateLimiter 'test' does not permit further calls" {
			t.Errorf("Expected error message 'RateLimiter 'test' does not permit further calls', but got '%s'", fullErr.Error())
		}
	}
}

func TestDecorateCover(t *testing.T) {
	rl := ratelimiter.NewRateLimiter("test")
	ratelimiter.DecorateRun(rl, func() {})
	ratelimiter.DecorateGet(rl, func() interface{} { return nil })
	ratelimiter.DecorateAccept(rl, func(interface{}) {})
	ratelimiter.DecorateApply(rl, func(_ interface{}) interface{} { return nil })
}
