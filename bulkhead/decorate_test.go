package bulkhead_test

import (
	"errors"
	"github.com/CharLemAznable/resilience4go/bulkhead"
	"github.com/stretchr/testify/assert"
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
	decoratedFn := bulkhead.DecorateRunnable(bh, fn)

	err1 := make(chan error, 1)
	err2 := make(chan error, 1)

	go func() {
		err1 <- decoratedFn()
	}()
	time.Sleep(time.Second * 1)
	go func() {
		err2 <- decoratedFn()
	}()

	assert.Equal(t, "error", (<-err1).Error())
	assert.Equal(t, "Bulkhead 'test' is full and does not permit further calls", (<-err2).Error())
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
	decoratedFn := bulkhead.DecorateSupplier(bh, fn)

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

	assert.Equal(t, "error", <-res1)
	assert.Equal(t, "error", (<-err1).Error())
	assert.Equal(t, "", <-res2)
	assert.Equal(t, "Bulkhead 'test' is full and does not permit further calls", (<-err2).Error())
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
	decoratedFn := bulkhead.DecorateConsumer(bh, fn)

	err1 := make(chan error, 1)
	err2 := make(chan error, 1)

	go func() {
		err1 <- decoratedFn("error")
	}()
	time.Sleep(time.Second * 1)
	go func() {
		err2 <- decoratedFn("failed")
	}()

	assert.Equal(t, "error", (<-err1).Error())
	assert.Equal(t, "Bulkhead 'test' is full and does not permit further calls", (<-err2).Error())
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
	decoratedFn := bulkhead.DecorateFunction(bh, fn)

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

	assert.Equal(t, "error", <-res1)
	assert.Equal(t, "error", (<-err1).Error())
	assert.Equal(t, "", <-res2)
	assert.Equal(t, "Bulkhead 'test' is full and does not permit further calls", (<-err2).Error())
}
