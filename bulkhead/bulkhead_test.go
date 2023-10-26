package bulkhead_test

import (
	"errors"
	"fmt"
	"github.com/CharLemAznable/resilience4go/bulkhead"
	"github.com/stretchr/testify/assert"
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
	eventProcessor := bh.EventProcessor()
	permitted := 0
	rejected := 0
	finished := 0
	eventProcessor.OnPermitted(func(event bulkhead.Event) {
		assert.Equal(t, bulkhead.PERMITTED, event.EventType())
		assert.Equal(t,
			fmt.Sprintf("%v: Bulkhead '%s' permitted a call.",
				event.CreationTime(), event.BulkheadName()),
			fmt.Sprintf("%v", event))
		permitted++
	})
	eventProcessor.OnRejected(func(event bulkhead.Event) {
		assert.Equal(t, bulkhead.REJECTED, event.EventType())
		assert.Equal(t,
			fmt.Sprintf("%v: Bulkhead '%s' rejected a call.",
				event.CreationTime(), event.BulkheadName()),
			fmt.Sprintf("%v", event))
		rejected++
	})
	eventProcessor.OnFinished(func(event bulkhead.Event) {
		assert.Equal(t, bulkhead.FINISHED, event.EventType())
		assert.Equal(t,
			fmt.Sprintf("%v: Bulkhead '%s' has finished a call.",
				event.CreationTime(), event.BulkheadName()),
			fmt.Sprintf("%v", event))
		finished++
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
	assert.Equal(t, 1, permitted)
	assert.Equal(t, 1, rejected)
	assert.Equal(t, 1, finished)
}
