package circuitbreaker_test

import (
	"errors"
	"github.com/CharLemAznable/resilience4go/circuitbreaker"
	"github.com/stretchr/testify/assert"
	"strconv"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestCircuitBreaker(t *testing.T) {
	breaker := circuitbreaker.NewCircuitBreaker("test",
		circuitbreaker.WithSlidingWindow(circuitbreaker.CountBased, 10, 10),
		circuitbreaker.WithAutomaticTransitionFromOpenToHalfOpenEnabled(true),
		circuitbreaker.WithWaitIntervalFunctionInOpenState(func(_ int64) time.Duration {
			return time.Second * 5
		}),
		circuitbreaker.WithPermittedNumberOfCallsInHalfOpenState(2))

	// 创建一个可运行的函数
	var count atomic.Int64
	fn := func() (string, error) {
		i := count.Add(1)
		str := strconv.FormatInt(i, 10)
		if i%2 == 0 {
			return str, nil
		}
		return "", errors.New(str)
	}
	// 调用DecorateRunnable函数
	decoratedFn := circuitbreaker.DecorateSupplier(breaker, fn)

	var wg sync.WaitGroup
	// 启动多个协程
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			_, _ = decoratedFn()
		}(i)
	}
	// 等待所有协程执行完毕
	wg.Wait()

	metrics := breaker.Metrics()
	assert.Equal(t, float64(50), metrics.FailureRate())
	assert.Equal(t, float64(0), metrics.SlowCallRate())
	assert.Equal(t, int64(10), metrics.NumberOfCalls())
	assert.Equal(t, int64(5), metrics.NumberOfSuccessfulCalls())
	assert.Equal(t, int64(5), metrics.NumberOfFailedCalls())
	assert.Equal(t, int64(0), metrics.NumberOfSlowCalls())
	assert.Equal(t, int64(0), metrics.NumberOfSlowSuccessfulCalls())
	assert.Equal(t, int64(0), metrics.NumberOfSlowFailedCalls())

	_, err := decoratedFn()
	e, ok := err.(*circuitbreaker.NotPermittedError)
	assert.True(t, ok)
	assert.Equal(t, "CircuitBreaker 'test' is OPEN and does not permit further calls", e.Error())

	metrics = breaker.Metrics()
	assert.Equal(t, int64(1), metrics.NumberOfNotPermittedCalls())

	time.Sleep(time.Second * 5)

	// HalfOpen
	_, err = decoratedFn()
	assert.Error(t, err)
	_, err = decoratedFn()
	assert.NoError(t, err)
	// Open again
	_, err = decoratedFn()
	e, ok = err.(*circuitbreaker.NotPermittedError)
	assert.True(t, ok)
	assert.Equal(t, "CircuitBreaker 'test' is OPEN and does not permit further calls", e.Error())

	time.Sleep(time.Second * 5)

	// HalfOpen
	count.Add(1)
	_, err = decoratedFn()
	assert.NoError(t, err)
	count.Add(1)
	_, err = decoratedFn()
	assert.NoError(t, err)
	// Closed
	count.Add(1)
	_, err = decoratedFn()
	assert.NoError(t, err)
}

func TestCircuitBreakerSlow(t *testing.T) {
	breaker := circuitbreaker.NewCircuitBreaker("test",
		circuitbreaker.WithSlidingWindow(circuitbreaker.TimeBased, 10, 10),
		circuitbreaker.WithSlowCallDurationThreshold(time.Second),
		circuitbreaker.WithWaitIntervalFunctionInOpenState(nil),
		circuitbreaker.WithPermittedNumberOfCallsInHalfOpenState(2),
		circuitbreaker.WithMaxWaitDurationInHalfOpenState(time.Second*5))

	// 创建一个可运行的函数
	fn := func(str string) (string, error) {
		time.Sleep(time.Second * 2)
		return str, nil
	}
	// 调用DecorateRunnable函数
	decoratedFn := circuitbreaker.DecorateFunction(breaker, fn)

	var wg sync.WaitGroup
	// 启动多个协程
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			_, _ = decoratedFn("test")
		}(i)
	}
	// 等待所有协程执行完毕
	wg.Wait()

	metrics := breaker.Metrics()
	assert.Equal(t, float64(0), metrics.FailureRate())
	assert.Equal(t, float64(100), metrics.SlowCallRate())
	assert.Equal(t, int64(10), metrics.NumberOfCalls())
	assert.Equal(t, int64(10), metrics.NumberOfSuccessfulCalls())
	assert.Equal(t, int64(0), metrics.NumberOfFailedCalls())
	assert.Equal(t, int64(10), metrics.NumberOfSlowCalls())
	assert.Equal(t, int64(10), metrics.NumberOfSlowSuccessfulCalls())
	assert.Equal(t, int64(0), metrics.NumberOfSlowFailedCalls())

	_, err := decoratedFn("test")
	e, ok := err.(*circuitbreaker.NotPermittedError)
	assert.True(t, ok)
	assert.Equal(t, "CircuitBreaker 'test' is OPEN and does not permit further calls", e.Error())

	metrics = breaker.Metrics()
	assert.Equal(t, int64(1), metrics.NumberOfNotPermittedCalls())

	_ = breaker.TransitionToHalfOpenState()

	_, err = decoratedFn("test")
	assert.NoError(t, err)

	time.Sleep(time.Second * 6)

	_, err = decoratedFn("test")
	e, ok = err.(*circuitbreaker.NotPermittedError)
	assert.True(t, ok)
	assert.Equal(t, "CircuitBreaker 'test' is OPEN and does not permit further calls", e.Error())

	_ = breaker.TransitionToHalfOpenState()

	// 启动多个协程
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			_, _ = decoratedFn("test")
		}(i)
	}
	// 等待所有协程执行完毕
	wg.Wait()

	_, err = decoratedFn("test")
	e, ok = err.(*circuitbreaker.NotPermittedError)
	assert.True(t, ok)
	assert.Equal(t, "CircuitBreaker 'test' is OPEN and does not permit further calls", e.Error())
}

func TestCircuitBreakerHalfOpenError(t *testing.T) {
	breaker := circuitbreaker.NewCircuitBreaker("halfOpenError")
	err := breaker.TransitionToHalfOpenState()
	expected := "CircuitBreaker 'halfOpenError' tried an illegal state transition from CLOSED to HALF_OPEN"
	assert.Equal(t, expected, err.Error())
}

func TestCircuitBreakerDisabled(t *testing.T) {
	breaker := circuitbreaker.NewCircuitBreaker("disabled",
		circuitbreaker.WithRecordResultPredicate(nil))
	_ = breaker.TransitionToDisabled()

	// 创建一个可运行的函数
	fn := func() error {
		panic("error")
	}
	// 调用DecorateRunnable函数
	decoratedFn := circuitbreaker.DecorateRunnable(breaker, fn)

	var wg sync.WaitGroup
	var count atomic.Int64
	// 启动多个协程
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			assert.PanicsWithValue(t, "error", func() {
				_ = decoratedFn()
			})
			count.Add(1)
		}(i)
	}
	// 等待所有协程执行完毕
	wg.Wait()
	assert.Equal(t, int64(100), count.Load())

	err := circuitbreaker.DecorateRunnable(breaker, func() error {
		return nil
	})()
	assert.NoError(t, err)
}

func TestCircuitBreakerForcedOpen(t *testing.T) {
	breaker := circuitbreaker.NewCircuitBreaker("forcedOpen")
	_ = breaker.TransitionToForcedOpen()

	// 创建一个可运行的函数
	fn := func(str string) error {
		return nil
	}
	// 调用DecorateRunnable函数
	decoratedFn := circuitbreaker.DecorateConsumer(breaker, fn)

	var wg sync.WaitGroup
	var count atomic.Int64
	// 启动多个协程
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			err := decoratedFn("test")
			e, ok := err.(*circuitbreaker.NotPermittedError)
			assert.True(t, ok)
			assert.Equal(t, "CircuitBreaker 'forcedOpen' is FORCED_OPEN and does not permit further calls", e.Error())
			count.Add(1)
		}(i)
	}
	// 等待所有协程执行完毕
	wg.Wait()
	assert.Equal(t, int64(100), count.Load())

	_ = breaker.TransitionToClosedState()
	err := decoratedFn("test")
	assert.NoError(t, err)
}
