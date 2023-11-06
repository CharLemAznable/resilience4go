package timelimiter_test

import (
	"errors"
	"github.com/CharLemAznable/resilience4go/timelimiter"
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

	func() {
		defer func() {
			if r := recover(); r != nil {
				if r != "panic error" {
					t.Errorf("Expected panic error 'panic error', but got '%v'", r)
				}
			}
		}()
		_ = decoratedFn()
	}()
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

	if ret != "error" {
		t.Errorf("Expected result 'error', but got '%s'", ret)
	}
	if err.Error() != "error" {
		t.Errorf("Expected error 'error', but got '%s'", err.Error())
	}
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
	if !ok {
		t.Errorf("Expected error type *timelimiter.TimeoutError, but got '%T'", err)
	} else {
		if timeout.Error() != "TimeLimiter 'test' recorded a timeout exception." {
			t.Errorf("Expected error message 'TimeLimiter 'test' recorded a timeout exception.', but got '%s'", timeout.Error())
		}
	}
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

	if ret != "" {
		t.Errorf("Expected result '', but got '%s'", ret)
	}
	timeout, ok := err.(*timelimiter.TimeoutError)
	if !ok {
		t.Errorf("Expected error type *timelimiter.TimeoutError, but got '%T'", err)
	} else {
		if timeout.Error() != "TimeLimiter 'test' recorded a timeout exception." {
			t.Errorf("Expected error message 'TimeLimiter 'test' recorded a timeout exception.', but got '%s'", timeout.Error())
		}
	}
}
