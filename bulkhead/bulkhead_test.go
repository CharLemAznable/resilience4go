package bulkhead_test

import (
	"errors"
	"fmt"
	"github.com/CharLemAznable/resilience4go/bulkhead"
	"sync/atomic"
	"testing"
	"time"
)

func TestBulkheadPublishEvents(t *testing.T) {
	// 创建一个Bulkhead的mock对象
	bh := bulkhead.NewBulkhead("test",
		bulkhead.WithMaxConcurrentCalls(1),
		bulkhead.WithMaxWaitDuration(time.Second*1))
	if bh.Name() != "test" {
		t.Errorf("Expected bulkhead name 'test', but got '%s'", bh.Name())
	}
	eventListener := bh.EventListener()
	permitted := atomic.Int64{}
	rejected := atomic.Int64{}
	finished := atomic.Int64{}
	onPermitted := func(event bulkhead.PermittedEvent) {
		if event.EventType() != bulkhead.PERMITTED {
			t.Errorf("Expected event type PERMITTED, but got '%s'", event.EventType())
		}
		expectedMsg := fmt.Sprintf("%v: Bulkhead '%s' permitted a call.", event.CreationTime(), event.BulkheadName())
		if event.String() != expectedMsg {
			t.Errorf("Expected event message '%s', but got '%s'", expectedMsg, event)
		}
		permitted.Add(1)
	}
	onRejected := func(event bulkhead.RejectedEvent) {
		if event.EventType() != bulkhead.REJECTED {
			t.Errorf("Expected event type REJECTED, but got '%s'", event.EventType())
		}
		expectedMsg := fmt.Sprintf("%v: Bulkhead '%s' rejected a call.", event.CreationTime(), event.BulkheadName())
		if event.String() != expectedMsg {
			t.Errorf("Expected event message '%s', but got '%s'", expectedMsg, event)
		}
		rejected.Add(1)
	}
	onFinished := func(event bulkhead.FinishedEvent) {
		if event.EventType() != bulkhead.FINISHED {
			t.Errorf("Expected event type FINISHED, but got '%s'", event.EventType())
		}
		expectedMsg := fmt.Sprintf("%v: Bulkhead '%s' has finished a call.", event.CreationTime(), event.BulkheadName())
		if event.String() != expectedMsg {
			t.Errorf("Expected event message '%s', but got '%s'", expectedMsg, event)
		}
		finished.Add(1)
	}
	eventListener.OnPermittedFunc(onPermitted).OnRejectedFunc(onRejected).OnFinishedFunc(onFinished)

	// 创建一个可运行的函数
	fn := func() error {
		if bh.Metrics().MaxAllowedConcurrentCalls() != 1 {
			t.Errorf("Expected MaxAllowedConcurrentCalls is 1, but got '%d'", bh.Metrics().MaxAllowedConcurrentCalls())
		}
		if bh.Metrics().AvailableConcurrentCalls() != 0 {
			t.Errorf("Expected AvailableConcurrentCalls is 0, but got '%d'", bh.Metrics().AvailableConcurrentCalls())
		}
		time.Sleep(time.Second * 3)
		return errors.New("error")
	}

	// 调用DecorateRunnable函数
	decoratedFn := bulkhead.DecorateRunnable(bh, fn)

	go func() {
		_ = decoratedFn()
	}()
	time.Sleep(time.Second * 1)
	go func() {
		_ = decoratedFn()
	}()

	time.Sleep(time.Second * 5)
	if permitted.Load() != 1 {
		t.Errorf("Expected 1 permitted call, but got '%d'", permitted.Load())
	}
	if rejected.Load() != 1 {
		t.Errorf("Expected 1 rejected call, but got '%d'", rejected.Load())
	}
	if finished.Load() != 1 {
		t.Errorf("Expected 1 finished call, but got '%d'", finished.Load())
	}
	eventListener.DismissPermittedFunc(onPermitted).DismissRejectedFunc(onRejected).DismissFinishedFunc(onFinished)
}
