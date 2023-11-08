package utils_test

import (
	"github.com/CharLemAznable/resilience4go/utils"
	"reflect"
	"testing"
)

func TestRemoveElementByValue(t *testing.T) {
	slice := []int{1, 2, 3, 3, 4, 5, 3}
	slices := utils.NewSlices[int](func(a, b int) bool {
		return a == b
	})
	result := slices.RemoveElementByValue(slice, 3)
	expected := []int{1, 2, 4, 5}
	if !slicesEqual(result, expected) {
		t.Errorf("Expected %v, but got %v", expected, result)
	}
	result = slices.RemoveElementByValue(slice, 0)
	expected = []int{1, 2, 3, 3, 4, 5, 3}
	if !slicesEqual(result, expected) {
		t.Errorf("Expected %v, but got %v", expected, result)
	}
	result = slices.AppendElementUnique(slice, 3)
	expected = []int{1, 2, 4, 5, 3}
	if !slicesEqual(result, expected) {
		t.Errorf("Expected %v, but got %v", expected, result)
	}
}

func slicesEqual(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}

type testFn func()

func TestRemoveElementByValue_Func(t *testing.T) {
	fn1 := func() {}
	fn2 := func() {}
	slice := []testFn{fn1, fn2, fn2, fn1}
	slices := utils.NewSlicesWithPointer[testFn]()
	result := slices.RemoveElementByValue(slice, fn2)
	expected := []testFn{fn1, fn1}
	if !slicesEqualFunc(result, expected) {
		t.Errorf("Expected %v, but got %v", expected, result)
	}
	fn3 := func() {}
	result = slices.RemoveElementByValue(slice, fn3)
	expected = []testFn{fn1, fn2, fn2, fn1}
	if !slicesEqualFunc(result, expected) {
		t.Errorf("Expected %v, but got %v", expected, result)
	}
	result = slices.AppendElementUnique(slice, fn2)
	expected = []testFn{fn1, fn1, fn2}
	if !slicesEqualFunc(result, expected) {
		t.Errorf("Expected %v, but got %v", expected, result)
	}
}

func slicesEqualFunc(a, b []testFn) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if reflect.ValueOf(v).Pointer() != reflect.ValueOf(b[i]).Pointer() {
			return false
		}
	}
	return true
}
