package ratelimiter_test

import (
	"fmt"
	"github.com/CharLemAznable/resilience4go/ratelimiter"
	"github.com/stretchr/testify/assert"
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
	assert.Equal(t, "test", rl.Name())
	eventListener := rl.EventListener()
	success := atomic.Int64{}
	failure := atomic.Int64{}
	eventListener.OnSuccess(func(event ratelimiter.Event) {
		assert.Equal(t, ratelimiter.SUCCESSFUL, event.EventType())
		assert.Equal(t,
			fmt.Sprintf("RateLimiterEvent{type=%s, rateLimiterName='%s', creationTime=%v}",
				event.EventType(), event.RateLimiterName(), event.CreationTime()),
			fmt.Sprintf("%v", event))
		success.Add(1)
	})
	eventListener.OnFailure(func(event ratelimiter.Event) {
		assert.Equal(t, ratelimiter.FAILED, event.EventType())
		assert.Equal(t,
			fmt.Sprintf("RateLimiterEvent{type=%s, rateLimiterName='%s', creationTime=%v}",
				event.EventType(), event.RateLimiterName(), event.CreationTime()),
			fmt.Sprintf("%v", event))
		failure.Add(1)
	})

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
	assert.Equal(t, int64(4), success.Load())
	assert.Equal(t, int64(1), failure.Load())
}
