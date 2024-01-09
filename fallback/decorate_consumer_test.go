package fallback_test

import (
	"errors"
	"github.com/CharLemAznable/resilience4go/fallback"
	"testing"
)

func TestDecorateConsumer(t *testing.T) {
	// Test case 1: fn returns nil error
	fn1 := func(str string) error {
		return nil
	}
	fallback1 := func(_ string) error {
		return errors.New("fallback error")
	}
	decoratedFn1 := fallback.DecorateCheckedAcceptDefault(fn1, fallback1)
	err1 := decoratedFn1("test")
	if err1 != nil {
		t.Errorf("Expected nil error, but got '%v'", err1)
	}

	// Test case 2: fn returns non-nil error
	fn2 := func(str string) error {
		return errors.New("original error")
	}
	fallback2 := func(_ string) error {
		return errors.New("fallback error")
	}
	decoratedFn2 := fallback.DecorateCheckedAcceptDefault(fn2, fallback2)
	err2 := decoratedFn2("test")
	if err2 == nil || err2.Error() != "fallback error" {
		t.Errorf("Expected error 'fallback error', but got '%v'", err2)
	}

	// Test case 3: fn returns TargetError
	fn3 := func(str string) error {
		return &TargetError{msg: "original error"}
	}
	fallback3 := func(_ string) error {
		return errors.New("fallback error")
	}
	decoratedFn3 := fallback.DecorateCheckedAcceptByType[string, *TargetError](fn3, fallback3)
	err3 := decoratedFn3("test")
	if err3 == nil || err3.Error() != "fallback error" {
		t.Errorf("Expected error 'fallback error', but got '%v'", err3)
	}

	// Test case 4: fn returns NonTargetError
	fn4 := func(str string) error {
		return &NonTargetError{msg: "original error"}
	}
	fallback4 := func(_ string, _ *TargetError) error {
		return errors.New("fallback error")
	}
	decoratedFn4 := fallback.DecorateCheckedAcceptWithFailure(fn4, fallback4)
	err4 := decoratedFn4("test")
	if err4 == nil || err4.Error() != "original error" {
		t.Errorf("Expected error 'original error', but got '%v'", err4)
	}

	// Test case 5: fn panic
	fn6 := func(str string) error {
		panic("original error")
	}
	fallback6 := func(_ string) error {
		return errors.New("fallback error")
	}
	decoratedFn6 := fallback.DecorateCheckedAcceptDefault(fn6, fallback6)
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
