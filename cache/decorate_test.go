package cache_test

import (
	"errors"
	"github.com/CharLemAznable/resilience4go/cache"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

type testKey struct {
	key string
}

type testValue struct {
	value string
}

func (v *testValue) String() string {
	return "testValue: [" + v.value + "]"
}

func TestDecorateFunction(t *testing.T) {
	ch := cache.NewCache[*testKey, *testValue]("test",
		cache.WithCapacity(10),
		cache.WithItemTTL(time.Second),
		cache.WithCacheResultPredicate(func(_ any, err error) bool {
			return err == nil
		}))

	var count atomic.Int64
	// fail with no error, max retries exceeded
	fn := func(key *testKey) (*testValue, error) {
		var err error
		if key.key == "nil" {
			count.Add(1)
			return nil, nil
		}
		if key.key == "ok" {
			err = errors.New("error")
		}
		return &testValue{key.key + randString(4)}, err
	}
	decoratedFn := cache.DecorateFunction(ch, fn)

	var wg sync.WaitGroup
	var ret1, ret2 *testValue
	wg.Add(1)
	go func() {
		ret1, _ = decoratedFn(&testKey{"notOK"})
		wg.Done()
	}()
	wg.Add(1)
	go func() {
		ret2, _ = decoratedFn(&testKey{"notOK"})
		wg.Done()
	}()
	wg.Wait()

	if ret1.value != ret2.value {
		t.Errorf("Expected return cached value, but got '%s' and '%s'", ret1.value, ret2.value)
	}

	time.Sleep(time.Second * 2)
	ret3, _ := decoratedFn(&testKey{"notOK"})
	if ret1.value == ret3.value {
		t.Errorf("Expected return new value, but got '%s'", ret3.value)
	}

	ret4, _ := decoratedFn(&testKey{"ok"})
	ret5, _ := decoratedFn(&testKey{"ok"})
	if ret4.value == ret5.value {
		t.Errorf("Expected return new value, but got '%s' and '%s'", ret4.value, ret5.value)
	}

	ret6, _ := decoratedFn(&testKey{"nil"})
	ret7, _ := decoratedFn(&testKey{"nil"})
	if ret6 != nil || ret7 != nil {
		t.Errorf("Expected return nil, but got '%v' and '%v'", ret6, ret6)
	}
	if count.Load() != 1 {
		t.Errorf("Expected count 1, but got %d", count.Load())
	}
}
