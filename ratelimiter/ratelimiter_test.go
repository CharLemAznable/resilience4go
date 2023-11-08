package ratelimiter_test

import (
	"fmt"
	"github.com/CharLemAznable/resilience4go/ratelimiter"
	"sync/atomic"
	"testing"
	"time"
)

func TestRateLimiterPublishEvents(t *testing.T) {
	// 创建一个TimeLimiter的mock对象
	rl := ratelimiter.NewRateLimiter("test",
		ratelimiter.WithTimeoutDuration(time.Second*2),
		ratelimiter.WithLimitRefreshPeriod(time.Second*2),
		ratelimiter.WithLimitForPeriod(2))
	if rl.Name() != "test" {
		t.Errorf("Expected rate limiter name 'test', but got '%s'", rl.Name())
	}
	metrics := rl.Metrics()
	if metrics.NumberOfWaitingThreads() != 0 {
		t.Errorf("Expected 0 waiting threads, but got '%d'", metrics.NumberOfWaitingThreads())
	}
	if metrics.AvailablePermissions() != 2 {
		t.Errorf("Expected 2 available permissions, but got '%d'", metrics.AvailablePermissions())
	}
	eventListener := rl.EventListener()
	success := atomic.Int64{}
	failure := atomic.Int64{}
	onSuccess := func(event ratelimiter.Event) {
		if event.EventType() != ratelimiter.SUCCESSFUL {
			t.Errorf("Expected event type SUCCESSFUL, but got '%s'", event.EventType())
		}
		expectedMsg := fmt.Sprintf("RateLimiterEvent{type=%s, rateLimiterName='%s', creationTime=%v}", event.EventType(), event.RateLimiterName(), event.CreationTime())
		if event.String() != expectedMsg {
			t.Errorf("Expected event message '%s', but got '%s'", expectedMsg, event)
		}
		success.Add(1)
	}
	onFailure := func(event ratelimiter.Event) {
		if event.EventType() != ratelimiter.FAILED {
			t.Errorf("Expected event type FAILED, but got '%s'", event.EventType())
		}
		expectedMsg := fmt.Sprintf("RateLimiterEvent{type=%s, rateLimiterName='%s', creationTime=%v}", event.EventType(), event.RateLimiterName(), event.CreationTime())
		if event.String() != expectedMsg {
			t.Errorf("Expected event message '%s', but got '%s'", expectedMsg, event)
		}
		failure.Add(1)
	}
	eventListener.OnSuccess(onSuccess).OnFailure(onFailure)
	if !eventListener.HasConsumer() {
		t.Error("Expected event listener has consumer, but not")
	}

	// 创建一个可运行的函数
	fn := func() error {
		time.Sleep(time.Millisecond * 500)
		return nil
	}
	// 调用DecorateRunnable函数
	decoratedFn := ratelimiter.DecorateRunnable(rl, fn)

	for i := 0; i < 5; i++ {
		go func() {
			_ = decoratedFn()
		}()
	}

	time.Sleep(time.Second * 5)
	if success.Load() != int64(4) {
		t.Errorf("Expected 4 successful calls, but got '%d'", success.Load())
	}
	if failure.Load() != int64(1) {
		t.Errorf("Expected 1 failure call, but got '%d'", failure.Load())
	}
	eventListener.Dismiss(onSuccess).Dismiss(onFailure)
	if eventListener.HasConsumer() {
		t.Error("Expected event listener has no consumer, but not")
	}
}
