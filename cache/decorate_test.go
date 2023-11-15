package cache_test

import (
	"github.com/CharLemAznable/resilience4go/cache"
	"sync"
	"testing"
	"time"
)

type testKey struct {
	key string
}

type testValue struct {
	value string
}

func TestDecorateFunction(t *testing.T) {
	ch := cache.NewCache[*testKey, *testValue]("test",
		cache.WithCapacity(10),
		cache.WithItemTTL(time.Second),
		cache.WithKeyToHash(nil))

	// fail with no error, max retries exceeded
	fn := func(key *testKey) (*testValue, error) {
		return &testValue{key.key + randString(4)}, nil
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
}
