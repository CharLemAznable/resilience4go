package retry_test

import (
	"errors"
	"github.com/CharLemAznable/resilience4go/retry"
	"sync/atomic"
	"testing"
)

func TestDecorateRunnable(t *testing.T) {
	rt := retry.NewRetry("test",
		retry.WithMaxAttempts(2),
		retry.WithFailAfterMaxAttempts(true))

	// panic
	fn := func() error {
		panic("panic")
	}
	decoratedFn := retry.DecorateRunnable(rt, fn)

	func() {
		defer func() {
			if r := recover(); r != nil {
				if r != "panic" {
					t.Errorf("Expected panic value 'panic', but got '%v'", r)
				}
			}
		}()
		_ = decoratedFn()
	}()
}

func TestDecorateSupplier(t *testing.T) {
	rt := retry.NewRetry("test",
		retry.WithMaxAttempts(2),
		retry.WithFailAfterMaxAttempts(true))

	// success directly
	fn := func() (int, error) {
		return 0, nil
	}
	decoratedFn := retry.DecorateSupplier(rt, fn)

	i, err := decoratedFn()
	if i != 0 {
		t.Errorf("Expected return value 0, but got '%d'", i)
	}
	if err != nil {
		t.Errorf("Expected nil error, but got '%v'", err)
	}
}

func TestDecorateConsumer(t *testing.T) {
	rt := retry.NewRetry("test",
		retry.WithMaxAttempts(2),
		retry.WithFailAfterMaxAttempts(true))

	// retry success
	var count atomic.Int64
	fn := func(val any) error {
		if count.Add(1) >= 2 {
			return nil
		}
		return errors.New("error")
	}
	decoratedFn := retry.DecorateConsumer(rt, fn)

	err := decoratedFn("test")
	if err != nil {
		t.Errorf("Expected nil error, but got '%v'", err)
	}
	if int(count.Load()) != 2 {
		t.Errorf("Expected count value 2, but got '%d'", int(count.Load()))
	}
}

func TestDecorateFunction(t *testing.T) {
	rt := retry.NewRetry("test",
		retry.WithMaxAttempts(2),
		retry.WithFailAfterMaxAttempts(true),
		retry.WithRecordResultPredicate(func(ret any, err error) bool {
			return ret.(string) != "ok" || err != nil
		}))

	// fail with no error, max retries exceeded
	fn := func(str string) (string, error) {
		return str, nil
	}
	decoratedFn := retry.DecorateFunction(rt, fn)

	ret, err := decoratedFn("notOK")
	if ret != "notOK" {
		t.Errorf("Expected return value 'notOK', but got '%s'", ret)
	}
	if err == nil {
		t.Error("Expected non-nil error")
	} else {
		expectedErr := "Retry 'test' has exhausted all attempts (2)"
		if err.Error() != expectedErr {
			t.Errorf("Expected error '%s', but got '%v'", expectedErr, err)
		}
	}
}
