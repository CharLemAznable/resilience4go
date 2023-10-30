package bulkhead_test

import (
	"errors"
	"fmt"
	"github.com/CharLemAznable/resilience4go/bulkhead"
	"github.com/stretchr/testify/assert"
	"sync/atomic"
	"testing"
	"time"
)

func TestBulkheadPublishEvents(t *testing.T) {
	// 创建一个可运行的函数
	fn := func() error {
		time.Sleep(time.Second * 3)
		return errors.New("error")
	}

	// 创建一个Bulkhead的mock对象
	bh := bulkhead.NewBulkhead("test",
		bulkhead.WithMaxConcurrentCalls(1),
		bulkhead.WithMaxWaitDuration(time.Second*1))
	assert.Equal(t, "test", bh.Name())
	eventListener := bh.EventListener()
	permitted := atomic.Int64{}
	rejected := atomic.Int64{}
	finished := atomic.Int64{}
	eventListener.OnPermitted(func(event bulkhead.Event) {
		assert.Equal(t, bulkhead.PERMITTED, event.EventType())
		assert.Equal(t,
			fmt.Sprintf("%v: Bulkhead '%s' permitted a call.",
				event.CreationTime(), event.BulkheadName()),
			fmt.Sprintf("%v", event))
		permitted.Add(1)
	})
	eventListener.OnRejected(func(event bulkhead.Event) {
		assert.Equal(t, bulkhead.REJECTED, event.EventType())
		assert.Equal(t,
			fmt.Sprintf("%v: Bulkhead '%s' rejected a call.",
				event.CreationTime(), event.BulkheadName()),
			fmt.Sprintf("%v", event))
		rejected.Add(1)
	})
	eventListener.OnFinished(func(event bulkhead.Event) {
		assert.Equal(t, bulkhead.FINISHED, event.EventType())
		assert.Equal(t,
			fmt.Sprintf("%v: Bulkhead '%s' has finished a call.",
				event.CreationTime(), event.BulkheadName()),
			fmt.Sprintf("%v", event))
		finished.Add(1)
	})

	// 调用DecorateRunnable函数
	decoratedFn := bulkhead.DecorateRunnable(bh, fn)

	go func() {
		decoratedFn.Run()
	}()
	time.Sleep(time.Second * 1)
	go func() {
		decoratedFn.Run()
	}()

	time.Sleep(time.Second * 5)
	assert.Equal(t, int64(1), permitted.Load())
	assert.Equal(t, int64(1), rejected.Load())
	assert.Equal(t, int64(1), finished.Load())
}
