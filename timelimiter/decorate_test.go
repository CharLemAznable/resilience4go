package timelimiter_test

import (
	"errors"
	"github.com/CharLemAznable/resilience4go/timelimiter"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestDecorateRunnable(t *testing.T) {
	// 创建一个可运行的函数
	fn := func() error {
		time.Sleep(time.Millisecond * 500)
		panic("panic error")
	}

	// 创建一个TimeLimiter的mock对象
	tl := timelimiter.NewTimeLimiter("test",
		timelimiter.WithTimeoutDuration(time.Second*1))

	// 调用DecorateRunnable函数
	decoratedFn := timelimiter.DecorateRunnable(tl, fn)

	assert.PanicsWithValue(t, "panic error", func() {
		_ = decoratedFn()
	})
}

func TestDecorateSupplier(t *testing.T) {
	// 创建一个可运行的函数
	fn := func() (string, error) {
		time.Sleep(time.Millisecond * 500)
		return "error", errors.New("error")
	}

	// 创建一个TimeLimiter的mock对象
	tl := timelimiter.NewTimeLimiter("test",
		timelimiter.WithTimeoutDuration(time.Second*1))

	// 调用DecorateRunnable函数
	decoratedFn := timelimiter.DecorateSupplier(tl, fn)

	ret, err := decoratedFn()

	assert.Equal(t, "error", ret)
	assert.Equal(t, "error", err.Error())
}

func TestDecorateConsumer(t *testing.T) {
	// 创建一个可运行的函数
	fn := func(str string) error {
		time.Sleep(time.Second * 2)
		return errors.New(str)
	}

	// 创建一个TimeLimiter的mock对象
	tl := timelimiter.NewTimeLimiter("test",
		timelimiter.WithTimeoutDuration(time.Second*1))

	// 调用DecorateRunnable函数
	decoratedFn := timelimiter.DecorateConsumer(tl, fn)

	err := decoratedFn("error")
	timeout, ok := err.(*timelimiter.TimeoutError)
	assert.True(t, ok)
	assert.Equal(t, "TimeLimiter 'test' recorded a timeout exception.", timeout.Error())
}

func TestDecorateFunction(t *testing.T) {
	// 创建一个可运行的函数
	fn := func(str string) (string, error) {
		time.Sleep(time.Second * 2)
		return str, errors.New(str)
	}

	// 创建一个TimeLimiter的mock对象
	tl := timelimiter.NewTimeLimiter("test",
		timelimiter.WithTimeoutDuration(time.Second*1))

	// 调用DecorateRunnable函数
	decoratedFn := timelimiter.DecorateFunction(tl, fn)

	ret, err := decoratedFn("error")

	assert.Equal(t, "", ret)
	timeout, ok := err.(*timelimiter.TimeoutError)
	assert.True(t, ok)
	assert.Equal(t, "TimeLimiter 'test' recorded a timeout exception.", timeout.Error())
}
