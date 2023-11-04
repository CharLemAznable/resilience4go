package retry_test

import (
	"errors"
	"github.com/CharLemAznable/resilience4go/retry"
	"github.com/stretchr/testify/assert"
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

	assert.PanicsWithValue(t, "panic", func() {
		_ = decoratedFn()
	})
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
	assert.Equal(t, 0, i)
	assert.NoError(t, err)
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
	assert.NoError(t, err)
	assert.Equal(t, 2, int(count.Load()))
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
	assert.Equal(t, "notOK", ret)
	assert.EqualError(t, err, "Retry 'test' has exhausted all attempts (2)")
}
