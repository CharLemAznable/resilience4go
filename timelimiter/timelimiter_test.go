package timelimiter_test

import (
	"errors"
	"fmt"
	"github.com/CharLemAznable/resilience4go/timelimiter"
	"github.com/stretchr/testify/assert"
	"sync/atomic"
	"testing"
	"time"
)

func TestTimeLimiterPublishEvents(t *testing.T) {
	// 创建一个TimeLimiter的mock对象
	tl := timelimiter.NewTimeLimiter("test",
		timelimiter.WithTimeoutDuration(time.Second*1))
	assert.Equal(t, "test", tl.Name())
	eventListener := tl.EventListener()
	success := atomic.Int64{}
	timeout := atomic.Int64{}
	failure := atomic.Int64{}
	panicMsg := "panic error"
	eventListener.OnSuccess(func(event timelimiter.Event) {
		assert.Equal(t, timelimiter.SUCCESS, event.EventType())
		assert.Equal(t,
			fmt.Sprintf("%v: TimeLimiter '%s' recorded a successful call.",
				event.CreationTime(), event.TimeLimiterName()),
			fmt.Sprintf("%v", event))
		success.Add(1)
	})
	eventListener.OnTimeout(func(event timelimiter.Event) {
		assert.Equal(t, timelimiter.TIMEOUT, event.EventType())
		assert.Equal(t,
			fmt.Sprintf("%v: TimeLimiter '%s' recorded a timeout exception.",
				event.CreationTime(), event.TimeLimiterName()),
			fmt.Sprintf("%v", event))
		timeout.Add(1)
	})
	eventListener.OnFailure(func(event timelimiter.Event) {
		assert.Equal(t, timelimiter.FAILURE, event.EventType())
		assert.Equal(t,
			fmt.Sprintf("%v: TimeLimiter '%s' recorded a failure call with panic: %v.",
				event.CreationTime(), event.TimeLimiterName(), panicMsg),
			fmt.Sprintf("%v", event))
		failure.Add(1)
	})

	// 创建一个可运行的函数
	fn := func() error {
		panic("panic error")
	}
	// 调用DecorateRunnable函数
	decoratedFn := timelimiter.DecorateRunnable(tl, fn)

	assert.PanicsWithValue(t, "panic error", func() {
		_ = decoratedFn()
	})

	// 创建一个可运行的函数
	fn = func() error {
		time.Sleep(time.Second * 2)
		return nil
	}
	// 调用DecorateRunnable函数
	decoratedFn = timelimiter.DecorateRunnable(tl, fn)

	err := decoratedFn()
	assert.Equal(t, "TimeLimiter 'test' recorded a timeout exception.", err.Error())

	// 创建一个可运行的函数
	fn = func() error {
		time.Sleep(time.Millisecond * 500)
		return errors.New("error")
	}
	// 调用DecorateRunnable函数
	decoratedFn = timelimiter.DecorateRunnable(tl, fn)

	err = decoratedFn()
	assert.Equal(t, "error", err.Error())

	time.Sleep(time.Second * 1)
	assert.Equal(t, int64(1), success.Load())
	assert.Equal(t, int64(1), timeout.Load())
	assert.Equal(t, int64(1), failure.Load())
}
