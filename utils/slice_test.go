package utils_test

import (
	"github.com/CharLemAznable/resilience4go/utils"
	"sync"
	"sync/atomic"
	"testing"
)

type testFn func()

func TestRemoveElementByValue_Func(t *testing.T) {
	fn1 := func() {}
	fn2 := func() {}
	slice := []testFn{fn1, fn2, fn2, fn1}
	result := utils.RemoveElementByValue(slice, fn2)
	expected := []testFn{fn1, fn1}
	if !slicesEqualFunc(result, expected) {
		t.Errorf("Expected %v, but got %v", expected, result)
	}
	fn3 := func() {}
	result = utils.RemoveElementByValue(slice, fn3)
	expected = []testFn{fn1, fn2, fn2, fn1}
	if !slicesEqualFunc(result, expected) {
		t.Errorf("Expected %v, but got %v", expected, result)
	}
	result = utils.AppendElementUnique(slice, fn2)
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
		if !utils.EqualsPointer(v, b[i]) {
			return false
		}
	}
	return true
}

func TestConsumeEvent(t *testing.T) {
	wg := &sync.WaitGroup{}
	wg.Add(1)
	var sum atomic.Int64
	fn1 := func(i int) {
		sum.Add(int64(i))
		wg.Done()
	}
	wg.Add(1)
	fn2 := func(i int) {
		sum.Add(int64(i))
		wg.Done()
	}
	slice := []func(int){fn1, fn2}
	utils.ConsumeEvent(slice, 2)
	wg.Wait()
	if sum.Load() != 4 {
		t.Errorf("Expected sum 4, but got %d", sum.Load())
	}
}
