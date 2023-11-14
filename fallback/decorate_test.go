package fallback_test

import (
	"errors"
	"github.com/CharLemAznable/resilience4go/fallback"
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
	decoratedFn1 := fallback.DecorateRunnableDefault(fn1, fallback1)
	err1 := decoratedFn1()
	if err1 != nil {
		t.Errorf("Expected nil error, but got '%v'", err1)
	}

	// Test case 2: fn returns non-nil error
	fn2 := func() error {
		return errors.New("original error")
	}
	fallback2 := func(err error) error {
		return errors.New("fallback error")
	}
	decoratedFn2 := fallback.DecorateRunnableDefault(fn2, fallback2)
	err2 := decoratedFn2()
	if err2 == nil || err2.Error() != "fallback error" {
		t.Errorf("Expected error 'fallback error', but got '%v'", err2)
	}

	// Test case 3: fn returns TargetError
	fn3 := func() error {
		return &TargetError{msg: "original error"}
	}
	fallback3 := func(err *TargetError) error {
		return errors.New("fallback error")
	}
	decoratedFn3 := fallback.DecorateRunnableDefault(fn3, fallback3)
	err3 := decoratedFn3()
	if err3 == nil || err3.Error() != "fallback error" {
		t.Errorf("Expected error 'fallback error', but got '%v'", err3)
	}

	// Test case 4: fn returns NonTargetError
	fn4 := func() error {
		return &NonTargetError{msg: "original error"}
	}
	fallback4 := func(err *TargetError) error {
		return errors.New("fallback error")
	}
	decoratedFn4 := fallback.DecorateRunnableDefault(fn4, fallback4)
	err4 := decoratedFn4()
	if err4 == nil || err4.Error() != "original error" {
		t.Errorf("Expected error 'original error', but got '%v'", err4)
	}

	// Test case 5: fn panic TargetError
	fn5 := func() error {
		panic(&TargetError{msg: "original error"})
	}
	fallback5 := func(err *TargetError) error {
		return errors.New("fallback error")
	}
	decoratedFn5 := fallback.DecorateRunnableDefault(fn5, fallback5)
	err5 := decoratedFn5()
	if err5 == nil || err5.Error() != "fallback error" {
		t.Errorf("Expected error 'fallback error', but got '%v'", err5)
	}

	// Test case 6: fn panic anything else
	fn6 := func() error {
		panic("original error")
	}
	fallback6 := func(err *TargetError) error {
		return errors.New("fallback error")
	}
	decoratedFn6 := fallback.DecorateRunnableDefault(fn6, fallback6)
	func() {
		defer func() {
			if r := recover(); r != nil {
				if r != "original error" {
					t.Errorf("Expected panic error 'panic error', but got '%v'", r)
				}
			}
		}()
		_ = decoratedFn6()
	}()
}

func TestDecorateSupplier(t *testing.T) {
	// Test case 1: fn returns nil error
	fn1 := func() (string, error) {
		return "ok", nil
	}
	fallback1 := func(_ string, err error) (string, error) {
		return "", errors.New("fallback error")
	}
	decoratedFn1 := fallback.DecorateSupplierDefault(fn1, fallback1)
	_, err1 := decoratedFn1()
	if err1 != nil {
		t.Errorf("Expected nil error, but got '%v'", err1)
	}

	// Test case 2: fn returns non-nil error
	fn2 := func() (string, error) {
		return "", errors.New("original error")
	}
	fallback2 := func(_ string, err error) (string, error) {
		return "", errors.New("fallback error")
	}
	decoratedFn2 := fallback.DecorateSupplierDefault(fn2, fallback2)
	_, err2 := decoratedFn2()
	if err2 == nil || err2.Error() != "fallback error" {
		t.Errorf("Expected error 'fallback error', but got '%v'", err2)
	}

	// Test case 3: fn returns TargetError
	fn3 := func() (string, error) {
		return "", &TargetError{msg: "original error"}
	}
	fallback3 := func(_ string, err *TargetError) (string, error) {
		return "", errors.New("fallback error")
	}
	decoratedFn3 := fallback.DecorateSupplierDefault(fn3, fallback3)
	_, err3 := decoratedFn3()
	if err3 == nil || err3.Error() != "fallback error" {
		t.Errorf("Expected error 'fallback error', but got '%v'", err3)
	}

	// Test case 4: fn returns NonTargetError
	fn4 := func() (string, error) {
		return "", &NonTargetError{msg: "original error"}
	}
	fallback4 := func(_ string, err *TargetError) (string, error) {
		return "", errors.New("fallback error")
	}
	decoratedFn4 := fallback.DecorateSupplierDefault(fn4, fallback4)
	_, err4 := decoratedFn4()
	if err4 == nil || err4.Error() != "original error" {
		t.Errorf("Expected error 'original error', but got '%v'", err4)
	}

	// Test case 5: fn panic TargetError
	fn5 := func() (string, error) {
		panic(&TargetError{msg: "original error"})
	}
	fallback5 := func(_ string, err *TargetError) (string, error) {
		return "", errors.New("fallback error")
	}
	decoratedFn5 := fallback.DecorateSupplierDefault(fn5, fallback5)
	_, err5 := decoratedFn5()
	if err5 == nil || err5.Error() != "fallback error" {
		t.Errorf("Expected error 'fallback error', but got '%v'", err5)
	}

	// Test case 6: fn panic anything else
	fn6 := func() (string, error) {
		panic("original error")
	}
	fallback6 := func(_ string, err *TargetError) (string, error) {
		return "", errors.New("fallback error")
	}
	decoratedFn6 := fallback.DecorateSupplierDefault(fn6, fallback6)
	func() {
		defer func() {
			if r := recover(); r != nil {
				if r != "original error" {
					t.Errorf("Expected panic error 'panic error', but got '%v'", r)
				}
			}
		}()
		_, _ = decoratedFn6()
	}()
}

func TestDecorateConsumer(t *testing.T) {
	// Test case 1: fn returns nil error
	fn1 := func(str string) error {
		return nil
	}
	fallback1 := func(_ string, err error) error {
		return errors.New("fallback error")
	}
	decoratedFn1 := fallback.DecorateConsumerDefault(fn1, fallback1)
	err1 := decoratedFn1("test")
	if err1 != nil {
		t.Errorf("Expected nil error, but got '%v'", err1)
	}

	// Test case 2: fn returns non-nil error
	fn2 := func(str string) error {
		return errors.New("original error")
	}
	fallback2 := func(_ string, err error) error {
		return errors.New("fallback error")
	}
	decoratedFn2 := fallback.DecorateConsumerDefault(fn2, fallback2)
	err2 := decoratedFn2("test")
	if err2 == nil || err2.Error() != "fallback error" {
		t.Errorf("Expected error 'fallback error', but got '%v'", err2)
	}

	// Test case 3: fn returns TargetError
	fn3 := func(str string) error {
		return &TargetError{msg: "original error"}
	}
	fallback3 := func(_ string, err *TargetError) error {
		return errors.New("fallback error")
	}
	decoratedFn3 := fallback.DecorateConsumerDefault(fn3, fallback3)
	err3 := decoratedFn3("test")
	if err3 == nil || err3.Error() != "fallback error" {
		t.Errorf("Expected error 'fallback error', but got '%v'", err3)
	}

	// Test case 4: fn returns NonTargetError
	fn4 := func(str string) error {
		return &NonTargetError{msg: "original error"}
	}
	fallback4 := func(_ string, err *TargetError) error {
		return errors.New("fallback error")
	}
	decoratedFn4 := fallback.DecorateConsumerDefault(fn4, fallback4)
	err4 := decoratedFn4("test")
	if err4 == nil || err4.Error() != "original error" {
		t.Errorf("Expected error 'original error', but got '%v'", err4)
	}

	// Test case 5: fn panic TargetError
	fn5 := func(str string) error {
		panic(&TargetError{msg: "original error"})
	}
	fallback5 := func(_ string, err *TargetError) error {
		return errors.New("fallback error")
	}
	decoratedFn5 := fallback.DecorateConsumerDefault(fn5, fallback5)
	err5 := decoratedFn5("test")
	if err5 == nil || err5.Error() != "fallback error" {
		t.Errorf("Expected error 'fallback error', but got '%v'", err5)
	}

	// Test case 6: fn panic anything else
	fn6 := func(str string) error {
		panic("original error")
	}
	fallback6 := func(_ string, err *TargetError) error {
		return errors.New("fallback error")
	}
	decoratedFn6 := fallback.DecorateConsumerDefault(fn6, fallback6)
	func() {
		defer func() {
			if r := recover(); r != nil {
				if r != "original error" {
					t.Errorf("Expected panic error 'panic error', but got '%v'", r)
				}
			}
		}()
		_ = decoratedFn6("test")
	}()
}

func TestDecorateFunction(t *testing.T) {
	// Test case 1: fn returns nil error
	fn1 := func(str string) (string, error) {
		return "ok", nil
	}
	fallback1 := func(_, _ string, err error) (string, error) {
		return "", errors.New("fallback error")
	}
	decoratedFn1 := fallback.DecorateFunctionDefault(fn1, fallback1)
	_, err1 := decoratedFn1("test")
	if err1 != nil {
		t.Errorf("Expected nil error, but got '%v'", err1)
	}

	// Test case 2: fn returns non-nil error
	fn2 := func(str string) (string, error) {
		return "", errors.New("original error")
	}
	fallback2 := func(_, _ string, err error) (string, error) {
		return "", errors.New("fallback error")
	}
	decoratedFn2 := fallback.DecorateFunctionDefault(fn2, fallback2)
	_, err2 := decoratedFn2("test")
	if err2 == nil || err2.Error() != "fallback error" {
		t.Errorf("Expected error 'fallback error', but got '%v'", err2)
	}

	// Test case 3: fn returns TargetError
	fn3 := func(str string) (string, error) {
		return "", &TargetError{msg: "original error"}
	}
	fallback3 := func(_, _ string, err *TargetError) (string, error) {
		return "", errors.New("fallback error")
	}
	decoratedFn3 := fallback.DecorateFunctionDefault(fn3, fallback3)
	_, err3 := decoratedFn3("test")
	if err3 == nil || err3.Error() != "fallback error" {
		t.Errorf("Expected error 'fallback error', but got '%v'", err3)
	}

	// Test case 4: fn returns NonTargetError
	fn4 := func(str string) (string, error) {
		return "", &NonTargetError{msg: "original error"}
	}
	fallback4 := func(_, _ string, err *TargetError) (string, error) {
		return "", errors.New("fallback error")
	}
	decoratedFn4 := fallback.DecorateFunctionDefault(fn4, fallback4)
	_, err4 := decoratedFn4("test")
	if err4 == nil || err4.Error() != "original error" {
		t.Errorf("Expected error 'original error', but got '%v'", err4)
	}

	// Test case 5: fn panic TargetError
	fn5 := func(str string) (string, error) {
		panic(&TargetError{msg: "original error"})
	}
	fallback5 := func(_, _ string, err *TargetError) (string, error) {
		return "", errors.New("fallback error")
	}
	decoratedFn5 := fallback.DecorateFunctionDefault(fn5, fallback5)
	_, err5 := decoratedFn5("test")
	if err5 == nil || err5.Error() != "fallback error" {
		t.Errorf("Expected error 'fallback error', but got '%v'", err5)
	}

	// Test case 6: fn panic anything else
	fn6 := func(str string) (string, error) {
		panic("original error")
	}
	fallback6 := func(_, _ string, err *TargetError) (string, error) {
		return "", errors.New("fallback error")
	}
	decoratedFn6 := fallback.DecorateFunctionDefault(fn6, fallback6)
	func() {
		defer func() {
			if r := recover(); r != nil {
				if r != "original error" {
					t.Errorf("Expected panic error 'panic error', but got '%v'", r)
				}
			}
		}()
		_, _ = decoratedFn6("test")
	}()
}
