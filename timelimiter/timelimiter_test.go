package timelimiter_test

import (
	"errors"
	"fmt"
	"github.com/CharLemAznable/resilience4go/timelimiter"
	"sync/atomic"
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
	success := atomic.Int64{}
	timeout := atomic.Int64{}
	failure := atomic.Int64{}
	panicMsg := "panic error"
	eventListener.OnSuccess(func(event timelimiter.Event) {
		if event.EventType() != timelimiter.SUCCESS {
			t.Errorf("Expected event type SUCCESS, but got '%s'", event.EventType())
		}
		expectedMsg := fmt.Sprintf("%v: TimeLimiter '%s' recorded a successful call.", event.CreationTime(), event.TimeLimiterName())
		if fmt.Sprintf("%v", event) != expectedMsg {
			t.Errorf("Expected event message '%s', but got '%s'", expectedMsg, fmt.Sprintf("%v", event))
		}
		success.Add(1)
	})
	eventListener.OnTimeout(func(event timelimiter.Event) {
		if event.EventType() != timelimiter.TIMEOUT {
			t.Errorf("Expected event type TIMEOUT, but got '%s'", event.EventType())
		}
		expectedMsg := fmt.Sprintf("%v: TimeLimiter '%s' recorded a timeout exception.", event.CreationTime(), event.TimeLimiterName())
		if fmt.Sprintf("%v", event) != expectedMsg {
			t.Errorf("Expected event message '%s', but got '%s'", expectedMsg, fmt.Sprintf("%v", event))
		}
		timeout.Add(1)
	})
	eventListener.OnFailure(func(event timelimiter.Event) {
		if event.EventType() != timelimiter.FAILURE {
			t.Errorf("Expected event type FAILURE, but got '%s'", event.EventType())
		}
		expectedMsg := fmt.Sprintf("%v: TimeLimiter '%s' recorded a failure call with panic: %v.", event.CreationTime(), event.TimeLimiterName(), panicMsg)
		if fmt.Sprintf("%v", event) != expectedMsg {
			t.Errorf("Expected event message '%s', but got '%s'", expectedMsg, fmt.Sprintf("%v", event))
		}
		failure.Add(1)
	})

	// 创建一个可运行的函数
	fn := func() error {
		panic("panic error")
	}
	// 调用DecorateRunnable函数
	decoratedFn := timelimiter.DecorateRunnable(tl, fn)

	func() {
		defer func() {
			if r := recover(); r != nil {
				if r != panicMsg {
					t.Errorf("Expected panic error '%s', but got '%v'", panicMsg, r)
				}
			}
		}()
		_ = decoratedFn()
	}()

	// 创建一个可运行的函数
	fn = func() error {
		time.Sleep(time.Second * 2)
		return nil
	}
	// 调用DecorateRunnable函数
	decoratedFn = timelimiter.DecorateRunnable(tl, fn)

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
	decoratedFn = timelimiter.DecorateRunnable(tl, fn)

	err = decoratedFn()
	if err.Error() != "error" {
		t.Errorf("Expected error 'error', but got '%s'", err.Error())
	}

	time.Sleep(time.Second * 1)
	if success.Load() != int64(1) {
		t.Errorf("Expected 1 success call, but got '%d'", success.Load())
	}
	if timeout.Load() != int64(1) {
		t.Errorf("Expected 1 timeout call, but got '%d'", timeout.Load())
	}
	if failure.Load() != int64(1) {
		t.Errorf("Expected 1 failure call, but got '%d'", failure.Load())
	}
}
