package timelimiter_test

import (
	"errors"
	"fmt"
	"github.com/CharLemAznable/resilience4go/timelimiter"
	"testing"
	"time"
)

func TestTimeLimiterPublishEvents(t *testing.T) {
	// 创建一个TimeLimiter的mock对象
	tl := timelimiter.NewTimeLimiter("test",
		timelimiter.WithTimeoutDuration(time.Second*1))
	if tl.Name() != "test" {
		t.Errorf("Expected time limiter name 'test', but got '%s'", tl.Name())
	}
	eventListener := tl.EventListener()
	onSuccess := func(event timelimiter.SuccessEvent) {
		if event.EventType() != timelimiter.SUCCESS {
			t.Errorf("Expected event type SUCCESS, but got '%s'", event.EventType())
		}
		expectedMsg := fmt.Sprintf("%v: TimeLimiter '%s' recorded a successful call.", event.CreationTime(), event.TimeLimiterName())
		if event.String() != expectedMsg {
			t.Errorf("Expected event message '%s', but got '%s'", expectedMsg, event)
		}
	}
	onError := func(event timelimiter.ErrorEvent) {
		if event.EventType() != timelimiter.ERROR {
			t.Errorf("Expected event type ERROR, but got '%s'", event.EventType())
		}
		expectedMsg := fmt.Sprintf("%v: TimeLimiter '%s' recorded a error call with error: %v.", event.CreationTime(), event.TimeLimiterName(), event.Error())
		if event.String() != expectedMsg {
			t.Errorf("Expected event message '%s', but got '%s'", expectedMsg, event)
		}
	}
	onTimeout := func(event timelimiter.TimeoutEvent) {
		if event.EventType() != timelimiter.TIMEOUT {
			t.Errorf("Expected event type TIMEOUT, but got '%s'", event.EventType())
		}
		expectedMsg := fmt.Sprintf("%v: TimeLimiter '%s' recorded a timeout exception.", event.CreationTime(), event.TimeLimiterName())
		if event.String() != expectedMsg {
			t.Errorf("Expected event message '%s', but got '%s'", expectedMsg, event)
		}
	}
	onPanic := func(event timelimiter.PanicEvent) {
		if event.EventType() != timelimiter.PANIC {
			t.Errorf("Expected event type PANIC, but got '%s'", event.EventType())
		}
		expectedMsg := fmt.Sprintf("%v: TimeLimiter '%s' recorded a failure call with panic: %v.", event.CreationTime(), event.TimeLimiterName(), event.Panic())
		if event.String() != expectedMsg {
			t.Errorf("Expected event message '%s', but got '%s'", expectedMsg, event)
		}
	}
	eventListener.OnSuccessFunc(onSuccess).OnErrorFunc(onError).OnTimeoutFunc(onTimeout).OnPanicFunc(onPanic)

	// 创建一个可运行的函数
	fn0 := func() {
		panic("panic error")
	}
	// 调用DecorateRunnable函数
	decoratedFn0 := timelimiter.DecorateRun(tl, fn0)

	func() {
		defer func() {
			if r := recover(); r != nil {
				if r != "panic error" {
					t.Errorf("Expected panic error 'panic error', but got '%v'", r)
				}
			}
		}()
		decoratedFn0()
	}()

	// 创建一个可运行的函数
	fn := func() error {
		time.Sleep(time.Second * 2)
		return nil
	}
	// 调用DecorateRunnable函数
	decoratedFn := timelimiter.DecorateCheckedRun(tl, fn)

	err := decoratedFn()
	if err.Error() != "TimeLimiter 'test' recorded a timeout exception." {
		t.Errorf("Expected error message 'TimeLimiter 'test' recorded a timeout exception.', but got '%s'", err.Error())
	}

	// 创建一个可运行的函数
	fn = func() error {
		time.Sleep(time.Millisecond * 500)
		return errors.New("error")
	}
	// 调用DecorateRunnable函数
	decoratedFn = timelimiter.DecorateCheckedRun(tl, fn)

	err = decoratedFn()
	if err.Error() != "error" {
		t.Errorf("Expected error 'error', but got '%s'", err.Error())
	}

	// 创建一个可运行的函数
	fn = func() error {
		time.Sleep(time.Millisecond * 500)
		return nil
	}
	// 调用DecorateRunnable函数
	decoratedFn = timelimiter.DecorateCheckedRun(tl, fn)

	err = decoratedFn()
	if err != nil {
		t.Errorf("Expected error is nil, but got '%s'", err.Error())
	}

	time.Sleep(time.Second * 1)
	metrics := tl.Metrics()
	if metrics.SuccessCount() != 1 {
		t.Errorf("Expected 1 success call, but got '%d'", metrics.SuccessCount())
	}
	if metrics.ErrorCount() != 1 {
		t.Errorf("Expected 1 error call, but got '%d'", metrics.ErrorCount())
	}
	if metrics.TimeoutCount() != 1 {
		t.Errorf("Expected 1 timeout call, but got '%d'", metrics.TimeoutCount())
	}
	if metrics.PanicCount() != 1 {
		t.Errorf("Expected 1 panic call, but got '%d'", metrics.PanicCount())
	}
	eventListener.DismissSuccessFunc(onSuccess).DismissErrorFunc(onError).DismissTimeoutFunc(onTimeout).DismissPanicFunc(onPanic)
}
