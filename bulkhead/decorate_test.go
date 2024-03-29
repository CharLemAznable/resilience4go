package bulkhead_test

import (
	"errors"
	"github.com/CharLemAznable/resilience4go/bulkhead"
	"testing"
	"time"
)

func TestDecorateRunnable(t *testing.T) {
	// 创建一个可运行的函数
	fn := func() error {
		time.Sleep(time.Second * 3)
		return errors.New("error")
	}

	// 创建一个Bulkhead的mock对象
	bh := bulkhead.NewBulkhead("test",
		bulkhead.WithMaxConcurrentCalls(1),
		bulkhead.WithMaxWaitDuration(time.Second*1))

	// 调用DecorateRunnable函数
	decoratedFn := bulkhead.DecorateCheckedRun(bh, fn)

	err1 := make(chan error, 1)
	err2 := make(chan error, 1)

	go func() {
		err1 <- decoratedFn()
	}()
	time.Sleep(time.Second * 1)
	go func() {
		err2 <- decoratedFn()
	}()

	if err := <-err1; err.Error() != "error" {
		t.Errorf("Expected error 'error', but got '%s'", err.Error())
	}
	err := <-err2
	fullErr, ok := err.(*bulkhead.FullError)
	if !ok {
		t.Errorf("Expected error type *bulkhead.FullError, but got '%T'", err)
	} else {
		if fullErr.Error() != "Bulkhead 'test' is full and does not permit further calls" {
			t.Errorf("Expected error message 'Bulkhead 'test' is full and does not permit further calls', but got '%s'", fullErr.Error())
		}
	}
}

func TestDecorateSupplier(t *testing.T) {
	// 创建一个可运行的函数
	fn := func() (string, error) {
		time.Sleep(time.Second * 3)
		return "error", errors.New("error")
	}

	// 创建一个Bulkhead的mock对象
	bh := bulkhead.NewBulkhead("test",
		bulkhead.WithMaxConcurrentCalls(1),
		bulkhead.WithMaxWaitDuration(time.Second*1))

	// 调用DecorateRunnable函数
	decoratedFn := bulkhead.DecorateCheckedGet(bh, fn)

	res1 := make(chan string, 1)
	err1 := make(chan error, 1)
	res2 := make(chan string, 1)
	err2 := make(chan error, 1)

	go func() {
		res, err := decoratedFn()
		res1 <- res
		err1 <- err
	}()
	time.Sleep(time.Second * 1)
	go func() {
		res, err := decoratedFn()
		res2 <- res
		err2 <- err
	}()

	if res := <-res1; res != "error" {
		t.Errorf("Expected result 'error', but got '%s'", res)
	}
	if err := <-err1; err.Error() != "error" {
		t.Errorf("Expected error 'error', but got '%s'", err.Error())
	}
	if res := <-res2; res != "" {
		t.Errorf("Expected result '', but got '%s'", res)
	}
	err := <-err2
	fullErr, ok := err.(*bulkhead.FullError)
	if !ok {
		t.Errorf("Expected error type *bulkhead.FullError, but got '%T'", err)
	} else {
		if fullErr.Error() != "Bulkhead 'test' is full and does not permit further calls" {
			t.Errorf("Expected error message 'Bulkhead 'test' is full and does not permit further calls', but got '%s'", fullErr.Error())
		}
	}
}

func TestDecorateConsumer(t *testing.T) {
	// 创建一个可运行的函数
	fn := func(str string) error {
		time.Sleep(time.Second * 3)
		return errors.New(str)
	}

	// 创建一个Bulkhead的mock对象
	bh := bulkhead.NewBulkhead("test",
		bulkhead.WithMaxConcurrentCalls(1),
		bulkhead.WithMaxWaitDuration(time.Second*1))

	// 调用DecorateRunnable函数
	decoratedFn := bulkhead.DecorateCheckedAccept(bh, fn)

	err1 := make(chan error, 1)
	err2 := make(chan error, 1)

	go func() {
		err1 <- decoratedFn("error")
	}()
	time.Sleep(time.Second * 1)
	go func() {
		err2 <- decoratedFn("failed")
	}()

	if err := <-err1; err.Error() != "error" {
		t.Errorf("Expected error 'error', but got '%s'", err.Error())
	}
	err := <-err2
	fullErr, ok := err.(*bulkhead.FullError)
	if !ok {
		t.Errorf("Expected error type *bulkhead.FullError, but got '%T'", err)
	} else {
		if fullErr.Error() != "Bulkhead 'test' is full and does not permit further calls" {
			t.Errorf("Expected error message 'Bulkhead 'test' is full and does not permit further calls', but got '%s'", fullErr.Error())
		}
	}
}

func TestDecorateFunction(t *testing.T) {
	// 创建一个可运行的函数
	fn := func(str string) (string, error) {
		time.Sleep(time.Second * 3)
		return str, errors.New(str)
	}

	// 创建一个Bulkhead的mock对象
	bh := bulkhead.NewBulkhead("test",
		bulkhead.WithMaxConcurrentCalls(1),
		bulkhead.WithMaxWaitDuration(time.Second*1))

	// 调用DecorateRunnable函数
	decoratedFn := bulkhead.DecorateCheckedApply(bh, fn)

	res1 := make(chan string, 1)
	err1 := make(chan error, 1)
	res2 := make(chan string, 1)
	err2 := make(chan error, 1)

	go func() {
		res, err := decoratedFn("error")
		res1 <- res
		err1 <- err
	}()
	time.Sleep(time.Second * 1)
	go func() {
		res, err := decoratedFn("failed")
		res2 <- res
		err2 <- err
	}()

	if res := <-res1; res != "error" {
		t.Errorf("Expected result 'error', but got '%s'", <-res1)
	}
	if err := <-err1; err.Error() != "error" {
		t.Errorf("Expected error 'error', but got '%s'", err.Error())
	}
	if res := <-res2; res != "" {
		t.Errorf("Expected result '', but got '%s'", res)
	}
	err := <-err2
	fullErr, ok := err.(*bulkhead.FullError)
	if !ok {
		t.Errorf("Expected error type *bulkhead.FullError, but got '%T'", err)
	} else {
		if fullErr.Error() != "Bulkhead 'test' is full and does not permit further calls" {
			t.Errorf("Expected error message 'Bulkhead 'test' is full and does not permit further calls', but got '%s'", fullErr.Error())
		}
	}
}

func TestDecorateCover(t *testing.T) {
	bh := bulkhead.NewBulkhead("test")
	bulkhead.DecorateRun(bh, func() {})
	bulkhead.DecorateGet(bh, func() interface{} { return nil })
	bulkhead.DecorateAccept(bh, func(interface{}) {})
	bulkhead.DecorateApply(bh, func(_ interface{}) interface{} { return nil })
}
