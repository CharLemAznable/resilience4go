package circuitbreaker_test

import (
	"errors"
	"fmt"
	"github.com/CharLemAznable/resilience4go/circuitbreaker"
	"strconv"
	"strings"
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
	state := breaker.State()
	if state != circuitbreaker.Closed {
		t.Errorf("Expected circuitbreaker state is CLOSED, but got %s", state)
	}
	listener := breaker.EventListener()
	listener.OnSuccess(func(event circuitbreaker.Event) {
		if event.EventType() != circuitbreaker.Success {
			t.Errorf("Expected event type Success, but got %s", event.EventType())
		}
		if !strings.HasPrefix(fmt.Sprintf("%v", event), fmt.Sprintf("%v: CircuitBreaker '%s'", event.CreationTime(), event.CircuitBreakerName())) {
			t.Errorf("Unexpected event message: %v", event)
		}
	})
	listener.OnError(func(event circuitbreaker.Event) {
		if event.EventType() != circuitbreaker.Error {
			t.Errorf("Expected event type Error, but got %s", event.EventType())
		}
		if !strings.HasPrefix(fmt.Sprintf("%v", event), fmt.Sprintf("%v: CircuitBreaker '%s'", event.CreationTime(), event.CircuitBreakerName())) {
			t.Errorf("Unexpected event message: %v", event)
		}
	})
	listener.OnNotPermitted(func(event circuitbreaker.Event) {
		if event.EventType() != circuitbreaker.NotPermitted {
			t.Errorf("Expected event type NotPermitted, but got %s", event.EventType())
		}
		if !strings.HasPrefix(fmt.Sprintf("%v", event), fmt.Sprintf("%v: CircuitBreaker '%s'", event.CreationTime(), event.CircuitBreakerName())) {
			t.Errorf("Unexpected event message: %v", event)
		}
	})
	listener.OnStateTransition(func(event circuitbreaker.Event) {
		if event.EventType() != circuitbreaker.StateTransition {
			t.Errorf("Expected event type StateTransition, but got %s", event.EventType())
		}
		if !strings.HasPrefix(fmt.Sprintf("%v", event), fmt.Sprintf("%v: CircuitBreaker '%s'", event.CreationTime(), event.CircuitBreakerName())) {
			t.Errorf("Unexpected event message: %v", event)
		}
	})
	listener.OnFailureRateExceeded(func(event circuitbreaker.Event) {
		if event.EventType() != circuitbreaker.FailureRateExceeded {
			t.Errorf("Expected event type FailureRateExceeded, but got %s", event.EventType())
		}
		if !strings.HasPrefix(fmt.Sprintf("%v", event), fmt.Sprintf("%v: CircuitBreaker '%s'", event.CreationTime(), event.CircuitBreakerName())) {
			t.Errorf("Unexpected event message: %v", event)
		}
	})
	listener.OnSlowCallRateExceeded(func(event circuitbreaker.Event) {
		t.Error("should not listen slow call rate exceeded event")
	})

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
	if metrics.FailureRate() != 50 {
		t.Errorf("Expected failure rate 50, but got %f", metrics.FailureRate())
	}
	if metrics.SlowCallRate() != 0 {
		t.Errorf("Expected slow call rate 0, but got %f", metrics.SlowCallRate())
	}
	if metrics.NumberOfCalls() != 10 {
		t.Errorf("Expected number of calls 10, but got %d", metrics.NumberOfCalls())
	}
	if metrics.NumberOfSuccessfulCalls() != 5 {
		t.Errorf("Expected number of successful calls 5, but got %d", metrics.NumberOfSuccessfulCalls())
	}
	if metrics.NumberOfFailedCalls() != 5 {
		t.Errorf("Expected number of failed calls 5, but got %d", metrics.NumberOfFailedCalls())
	}
	if metrics.NumberOfSlowCalls() != 0 {
		t.Errorf("Expected number of slow calls 0, but got %d", metrics.NumberOfSlowCalls())
	}
	if metrics.NumberOfSlowSuccessfulCalls() != 0 {
		t.Errorf("Expected number of slow successful calls 0, but got %d", metrics.NumberOfSlowSuccessfulCalls())
	}
	if metrics.NumberOfSlowFailedCalls() != 0 {
		t.Errorf("Expected number of slow failed calls 0, but got %d", metrics.NumberOfSlowFailedCalls())
	}

	_, err := decoratedFn()
	e, ok := err.(*circuitbreaker.NotPermittedError)
	if !ok {
		t.Errorf("Expected error type *circuitbreaker.NotPermittedError, but got %T", err)
	}
	if e.Error() != "CircuitBreaker 'test' is OPEN and does not permit further calls" {
		t.Errorf("Expected error message 'CircuitBreaker 'test' is OPEN and does not permit further calls', but got '%s'", e.Error())
	}
	state = breaker.State()
	if state != circuitbreaker.Open {
		t.Errorf("Expected circuitbreaker state is OPEN, but got %s", state)
	}
	metrics = breaker.Metrics()
	if metrics.NumberOfNotPermittedCalls() != 1 {
		t.Errorf("Expected number of not permitted calls 1, but got %d", metrics.NumberOfNotPermittedCalls())
	}

	time.Sleep(time.Second * 5)

	// HalfOpen
	_, err = decoratedFn()
	if err == nil {
		t.Error("Expected error, but got nil")
	}
	state = breaker.State()
	if state != circuitbreaker.HalfOpen {
		t.Errorf("Expected circuitbreaker state is HALF_OPEN, but got %s", state)
	}
	_, err = decoratedFn()
	if err != nil {
		t.Errorf("Expected nil error, but got %v", err)
	}
	// Open again
	_, err = decoratedFn()
	e, ok = err.(*circuitbreaker.NotPermittedError)
	if !ok {
		t.Errorf("Expected error type *circuitbreaker.NotPermittedError, but got %T", err)
	}
	if e.Error() != "CircuitBreaker 'test' is OPEN and does not permit further calls" {
		t.Errorf("Expected error message 'CircuitBreaker 'test' is OPEN and does not permit further calls', but got '%s'", e.Error())
	}
	state = breaker.State()
	if state != circuitbreaker.Open {
		t.Errorf("Expected circuitbreaker state is OPEN, but got %s", state)
	}

	time.Sleep(time.Second * 5)

	// HalfOpen
	count.Add(1)
	_, err = decoratedFn()
	if err != nil {
		t.Errorf("Expected nil error, but got %v", err)
	}
	state = breaker.State()
	if state != circuitbreaker.HalfOpen {
		t.Errorf("Expected circuitbreaker state is HALF_OPEN, but got %s", state)
	}
	count.Add(1)
	_, err = decoratedFn()
	if err != nil {
		t.Errorf("Expected nil error, but got %v", err)
	}
	// Closed
	count.Add(1)
	_, err = decoratedFn()
	if err != nil {
		t.Errorf("Expected nil error, but got %v", err)
	}
	state = breaker.State()
	if state != circuitbreaker.Closed {
		t.Errorf("Expected circuitbreaker state is CLOSED, but got %s", state)
	}

	time.Sleep(time.Second)
}

func TestCircuitBreakerSlow(t *testing.T) {
	breaker := circuitbreaker.NewCircuitBreaker("test",
		circuitbreaker.WithSlidingWindow(circuitbreaker.TimeBased, 10, 10),
		circuitbreaker.WithSlowCallDurationThreshold(time.Second),
		circuitbreaker.WithWaitIntervalFunctionInOpenState(nil),
		circuitbreaker.WithPermittedNumberOfCallsInHalfOpenState(2),
		circuitbreaker.WithMaxWaitDurationInHalfOpenState(time.Second*5))
	listener := breaker.EventListener()
	listener.OnSuccess(func(event circuitbreaker.Event) {
		if event.EventType() != circuitbreaker.Success {
			t.Errorf("Expected event type Success, but got %s", event.EventType())
		}
		if !strings.HasPrefix(fmt.Sprintf("%v", event), fmt.Sprintf("%v: CircuitBreaker '%s'", event.CreationTime(), event.CircuitBreakerName())) {
			t.Errorf("Unexpected event message: %v", event)
		}
	})
	listener.OnError(func(event circuitbreaker.Event) {
		t.Error("should not listen error event")
	})
	listener.OnNotPermitted(func(event circuitbreaker.Event) {
		if event.EventType() != circuitbreaker.NotPermitted {
			t.Errorf("Expected event type NotPermitted, but got %s", event.EventType())
		}
		if !strings.HasPrefix(fmt.Sprintf("%v", event), fmt.Sprintf("%v: CircuitBreaker '%s'", event.CreationTime(), event.CircuitBreakerName())) {
			t.Errorf("Unexpected event message: %v", event)
		}
	})
	listener.OnStateTransition(func(event circuitbreaker.Event) {
		if event.EventType() != circuitbreaker.StateTransition {
			t.Errorf("Expected event type StateTransition, but got %s", event.EventType())
		}
		if !strings.HasPrefix(fmt.Sprintf("%v", event), fmt.Sprintf("%v: CircuitBreaker '%s'", event.CreationTime(), event.CircuitBreakerName())) {
			t.Errorf("Unexpected event message: %v", event)
		}
	})
	listener.OnFailureRateExceeded(func(event circuitbreaker.Event) {
		t.Error("should not listen failure rate exceeded event")
	})
	listener.OnSlowCallRateExceeded(func(event circuitbreaker.Event) {
		if event.EventType() != circuitbreaker.SlowCallRateExceeded {
			t.Errorf("Expected event type SlowCallRateExceeded, but got %s", event.EventType())
		}
		if !strings.HasPrefix(fmt.Sprintf("%v", event), fmt.Sprintf("%v: CircuitBreaker '%s'", event.CreationTime(), event.CircuitBreakerName())) {
			t.Errorf("Unexpected event message: %v", event)
		}
	})

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
	if metrics.FailureRate() != 0 {
		t.Errorf("Expected failure rate 0, but got %f", metrics.FailureRate())
	}
	if metrics.SlowCallRate() != 100 {
		t.Errorf("Expected slow call rate 100, but got %f", metrics.SlowCallRate())
	}
	if metrics.NumberOfCalls() != 10 {
		t.Errorf("Expected number of calls 10, but got %d", metrics.NumberOfCalls())
	}
	if metrics.NumberOfSuccessfulCalls() != 10 {
		t.Errorf("Expected number of successful calls 10, but got %d", metrics.NumberOfSuccessfulCalls())
	}
	if metrics.NumberOfFailedCalls() != 0 {
		t.Errorf("Expected number of failed calls 0, but got %d", metrics.NumberOfFailedCalls())
	}
	if metrics.NumberOfSlowCalls() != 10 {
		t.Errorf("Expected number of slow calls 10, but got %d", metrics.NumberOfSlowCalls())
	}
	if metrics.NumberOfSlowSuccessfulCalls() != 10 {
		t.Errorf("Expected number of slow successful calls 10, but got %d", metrics.NumberOfSlowSuccessfulCalls())
	}
	if metrics.NumberOfSlowFailedCalls() != 0 {
		t.Errorf("Expected number of slow failed calls 0, but got %d", metrics.NumberOfSlowFailedCalls())
	}

	_, err := decoratedFn("test")
	e, ok := err.(*circuitbreaker.NotPermittedError)
	if !ok {
		t.Errorf("Expected error type *circuitbreaker.NotPermittedError, but got %T", err)
	}
	if e.Error() != "CircuitBreaker 'test' is OPEN and does not permit further calls" {
		t.Errorf("Expected error message 'CircuitBreaker 'test' is OPEN and does not permit further calls', but got '%s'", e.Error())
	}

	metrics = breaker.Metrics()
	if metrics.NumberOfNotPermittedCalls() != 1 {
		t.Errorf("Expected number of not permitted calls 1, but got %d", metrics.NumberOfNotPermittedCalls())
	}

	_ = breaker.TransitionToHalfOpenState()

	_, err = decoratedFn("test")
	if err != nil {
		t.Errorf("Expected nil error, but got %v", err)
	}

	time.Sleep(time.Second * 6)

	_, err = decoratedFn("test")
	e, ok = err.(*circuitbreaker.NotPermittedError)
	if !ok {
		t.Errorf("Expected error type *circuitbreaker.NotPermittedError, but got %T", err)
	}
	if e.Error() != "CircuitBreaker 'test' is OPEN and does not permit further calls" {
		t.Errorf("Expected error message 'CircuitBreaker 'test' is OPEN and does not permit further calls', but got '%s'", e.Error())
	}

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
	if !ok {
		t.Errorf("Expected error type *circuitbreaker.NotPermittedError, but got %T", err)
	}
	if e.Error() != "CircuitBreaker 'test' is OPEN and does not permit further calls" {
		t.Errorf("Expected error message 'CircuitBreaker 'test' is OPEN and does not permit further calls', but got '%s'", e.Error())
	}

	time.Sleep(time.Second)
}

func TestCircuitBreakerHalfOpenError(t *testing.T) {
	breaker := circuitbreaker.NewCircuitBreaker("halfOpenError")
	err := breaker.TransitionToHalfOpenState()
	expected := "CircuitBreaker 'halfOpenError' tried an illegal state transition from CLOSED to HALF_OPEN"
	if err.Error() != expected {
		t.Errorf("Expected error message '%s', but got '%s'", expected, err.Error())
	}
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
			func() {
				defer func() {
					if r := recover(); r != nil {
						if r != "error" {
							t.Errorf("Expected panic error 'error', but got '%v'", r)
						}
					}
				}()
				_ = decoratedFn()
			}()
			count.Add(1)
		}(i)
	}
	// 等待所有协程执行完毕
	wg.Wait()
	if count.Load() != 100 {
		t.Errorf("Expected count 100, but got %d", count.Load())
	}

	err := circuitbreaker.DecorateRunnable(breaker, func() error {
		return nil
	})()
	if err != nil {
		t.Errorf("Expected nil error, but got %v", err)
	}
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
			if !ok {
				t.Errorf("Expected error type *circuitbreaker.NotPermittedError, but got %T", err)
			}
			expected := "CircuitBreaker 'forcedOpen' is FORCED_OPEN and does not permit further calls"
			if e.Error() != expected {
				t.Errorf("Expected error message '%s', but got '%s'", expected, e.Error())
			}
			count.Add(1)
		}(i)
	}
	// 等待所有协程执行完毕
	wg.Wait()
	if count.Load() != 100 {
		t.Errorf("Expected count 100, but got %d", count.Load())
	}

	_ = breaker.TransitionToClosedState()
	err := decoratedFn("test")
	if err != nil {
		t.Errorf("Expected nil error, but got %v", err)
	}
}
