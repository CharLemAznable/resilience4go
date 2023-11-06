package fallback_test

import (
	"errors"
	"github.com/CharLemAznable/resilience4go/fallback"
	"github.com/stretchr/testify/assert"
	"testing"
)

type TargetError struct {
	msg string
}

func (e *TargetError) Error() string {
	return e.msg
}

type NonTargetError struct {
	msg string
}

func (e *NonTargetError) Error() string {
	return e.msg
}

func TestDecorateRunnable(t *testing.T) {
	// Test case 1: fn returns nil error
	fn1 := func() error {
		return nil
	}
	fallback1 := func(err error) error {
		return errors.New("fallback error")
	}
	decoratedFn1 := fallback.DecorateRunnable(fn1, fallback1)
	err1 := decoratedFn1()
	assert.NoError(t, err1)

	// Test case 2: fn returns non-nil error
	fn2 := func() error {
		return errors.New("original error")
	}
	fallback2 := func(err error) error {
		return errors.New("fallback error")
	}
	decoratedFn2 := fallback.DecorateRunnable(fn2, fallback2)
	err2 := decoratedFn2()
	assert.EqualError(t, err2, "fallback error")

	// Test case 3: fn returns TargetError
	fn3 := func() error {
		return &TargetError{msg: "original error"}
	}
	fallback3 := func(err *TargetError) error {
		return errors.New("fallback error")
	}
	decoratedFn3 := fallback.DecorateRunnable(fn3, fallback3)
	err3 := decoratedFn3()
	assert.EqualError(t, err3, "fallback error")

	// Test case 4: fn returns NonTargetError
	fn4 := func() error {
		return &NonTargetError{msg: "original error"}
	}
	fallback4 := func(err *TargetError) error {
		return errors.New("fallback error")
	}
	decoratedFn4 := fallback.DecorateRunnable(fn4, fallback4)
	err4 := decoratedFn4()
	assert.EqualError(t, err4, "original error")

	// Test case 5: fn panic TargetError
	fn5 := func() error {
		panic(&TargetError{msg: "original error"})
	}
	fallback5 := func(err *TargetError) error {
		return errors.New("fallback error")
	}
	decoratedFn5 := fallback.DecorateRunnable(fn5, fallback5)
	err5 := decoratedFn5()
	assert.EqualError(t, err5, "fallback error")

	// Test case 6: fn panic anything else
	fn6 := func() error {
		panic("original error")
	}
	fallback6 := func(err *TargetError) error {
		return errors.New("fallback error")
	}
	decoratedFn6 := fallback.DecorateRunnable(fn6, fallback6)
	assert.PanicsWithValue(t, "original error", func() {
		_ = decoratedFn6()
	})
}

func TestDecorateSupplier(t *testing.T) {
	// Test case 1: fn returns nil error
	fn1 := func() (string, error) {
		return "ok", nil
	}
	fallback1 := func(err error) (string, error) {
		return "", errors.New("fallback error")
	}
	decoratedFn1 := fallback.DecorateSupplier(fn1, fallback1)
	_, err1 := decoratedFn1()
	assert.NoError(t, err1)

	// Test case 2: fn returns non-nil error
	fn2 := func() (string, error) {
		return "", errors.New("original error")
	}
	fallback2 := func(err error) (string, error) {
		return "", errors.New("fallback error")
	}
	decoratedFn2 := fallback.DecorateSupplier(fn2, fallback2)
	_, err2 := decoratedFn2()
	assert.EqualError(t, err2, "fallback error")

	// Test case 3: fn returns TargetError
	fn3 := func() (string, error) {
		return "", &TargetError{msg: "original error"}
	}
	fallback3 := func(err *TargetError) (string, error) {
		return "", errors.New("fallback error")
	}
	decoratedFn3 := fallback.DecorateSupplier(fn3, fallback3)
	_, err3 := decoratedFn3()
	assert.EqualError(t, err3, "fallback error")

	// Test case 4: fn returns NonTargetError
	fn4 := func() (string, error) {
		return "", &NonTargetError{msg: "original error"}
	}
	fallback4 := func(err *TargetError) (string, error) {
		return "", errors.New("fallback error")
	}
	decoratedFn4 := fallback.DecorateSupplier(fn4, fallback4)
	_, err4 := decoratedFn4()
	assert.EqualError(t, err4, "original error")

	// Test case 5: fn panic TargetError
	fn5 := func() (string, error) {
		panic(&TargetError{msg: "original error"})
	}
	fallback5 := func(err *TargetError) (string, error) {
		return "", errors.New("fallback error")
	}
	decoratedFn5 := fallback.DecorateSupplier(fn5, fallback5)
	_, err5 := decoratedFn5()
	assert.EqualError(t, err5, "fallback error")

	// Test case 6: fn panic anything else
	fn6 := func() (string, error) {
		panic("original error")
	}
	fallback6 := func(err *TargetError) (string, error) {
		return "", errors.New("fallback error")
	}
	decoratedFn6 := fallback.DecorateSupplier(fn6, fallback6)
	assert.PanicsWithValue(t, "original error", func() {
		_, _ = decoratedFn6()
	})
}

func TestDecorateConsumer(t *testing.T) {
	// Test case 1: fn returns nil error
	fn1 := func(str string) error {
		return nil
	}
	fallback1 := func(err error) error {
		return errors.New("fallback error")
	}
	decoratedFn1 := fallback.DecorateConsumer(fn1, fallback1)
	err1 := decoratedFn1("test")
	assert.NoError(t, err1)

	// Test case 2: fn returns non-nil error
	fn2 := func(str string) error {
		return errors.New("original error")
	}
	fallback2 := func(err error) error {
		return errors.New("fallback error")
	}
	decoratedFn2 := fallback.DecorateConsumer(fn2, fallback2)
	err2 := decoratedFn2("test")
	assert.EqualError(t, err2, "fallback error")

	// Test case 3: fn returns TargetError
	fn3 := func(str string) error {
		return &TargetError{msg: "original error"}
	}
	fallback3 := func(err *TargetError) error {
		return errors.New("fallback error")
	}
	decoratedFn3 := fallback.DecorateConsumer(fn3, fallback3)
	err3 := decoratedFn3("test")
	assert.EqualError(t, err3, "fallback error")

	// Test case 4: fn returns NonTargetError
	fn4 := func(str string) error {
		return &NonTargetError{msg: "original error"}
	}
	fallback4 := func(err *TargetError) error {
		return errors.New("fallback error")
	}
	decoratedFn4 := fallback.DecorateConsumer(fn4, fallback4)
	err4 := decoratedFn4("test")
	assert.EqualError(t, err4, "original error")

	// Test case 5: fn panic TargetError
	fn5 := func(str string) error {
		panic(&TargetError{msg: "original error"})
	}
	fallback5 := func(err *TargetError) error {
		return errors.New("fallback error")
	}
	decoratedFn5 := fallback.DecorateConsumer(fn5, fallback5)
	err5 := decoratedFn5("test")
	assert.EqualError(t, err5, "fallback error")

	// Test case 6: fn panic anything else
	fn6 := func(str string) error {
		panic("original error")
	}
	fallback6 := func(err *TargetError) error {
		return errors.New("fallback error")
	}
	decoratedFn6 := fallback.DecorateConsumer(fn6, fallback6)
	assert.PanicsWithValue(t, "original error", func() {
		_ = decoratedFn6("test")
	})
}

func TestDecorateFunction(t *testing.T) {
	// Test case 1: fn returns nil error
	fn1 := func(str string) (string, error) {
		return "ok", nil
	}
	fallback1 := func(err error) (string, error) {
		return "", errors.New("fallback error")
	}
	decoratedFn1 := fallback.DecorateFunction(fn1, fallback1)
	_, err1 := decoratedFn1("test")
	assert.NoError(t, err1)

	// Test case 2: fn returns non-nil error
	fn2 := func(str string) (string, error) {
		return "", errors.New("original error")
	}
	fallback2 := func(err error) (string, error) {
		return "", errors.New("fallback error")
	}
	decoratedFn2 := fallback.DecorateFunction(fn2, fallback2)
	_, err2 := decoratedFn2("test")
	assert.EqualError(t, err2, "fallback error")

	// Test case 3: fn returns TargetError
	fn3 := func(str string) (string, error) {
		return "", &TargetError{msg: "original error"}
	}
	fallback3 := func(err *TargetError) (string, error) {
		return "", errors.New("fallback error")
	}
	decoratedFn3 := fallback.DecorateFunction(fn3, fallback3)
	_, err3 := decoratedFn3("test")
	assert.EqualError(t, err3, "fallback error")

	// Test case 4: fn returns NonTargetError
	fn4 := func(str string) (string, error) {
		return "", &NonTargetError{msg: "original error"}
	}
	fallback4 := func(err *TargetError) (string, error) {
		return "", errors.New("fallback error")
	}
	decoratedFn4 := fallback.DecorateFunction(fn4, fallback4)
	_, err4 := decoratedFn4("test")
	assert.EqualError(t, err4, "original error")

	// Test case 5: fn panic TargetError
	fn5 := func(str string) (string, error) {
		panic(&TargetError{msg: "original error"})
	}
	fallback5 := func(err *TargetError) (string, error) {
		return "", errors.New("fallback error")
	}
	decoratedFn5 := fallback.DecorateFunction(fn5, fallback5)
	_, err5 := decoratedFn5("test")
	assert.EqualError(t, err5, "fallback error")

	// Test case 6: fn panic anything else
	fn6 := func(str string) (string, error) {
		panic("original error")
	}
	fallback6 := func(err *TargetError) (string, error) {
		return "", errors.New("fallback error")
	}
	decoratedFn6 := fallback.DecorateFunction(fn6, fallback6)
	assert.PanicsWithValue(t, "original error", func() {
		_, _ = decoratedFn6("test")
	})
}
