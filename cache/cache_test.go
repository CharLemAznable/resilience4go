package cache_test

import (
	"fmt"
	"github.com/CharLemAznable/resilience4go/cache"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestCache(t *testing.T) {
	ch := cache.NewCache[string, string]("test",
		cache.WithCapacity(10),
		cache.WithItemTTL(time.Second),
		cache.WithKeyToHash(nil),
		cache.WithCacheResultPredicate(nil)).
		WithMarshalFn(func(s string) any {
			return []byte(s)
		}, func(v any) string {
			return string(v.([]byte))
		})
	if ch.Name() != "test" {
		t.Errorf("Expected cache name 'test', but got '%s'", ch.Name())
	}
	eventListener := ch.EventListener()
	var hits atomic.Int64
	var misses atomic.Int64
	onCacheHit := func(event cache.HitEvent) {
		if event.EventType() != cache.OnHit {
			t.Errorf("Expected event type CACHE_HIT, but got '%s'", event.EventType())
		}
		expectedMsg := fmt.Sprintf("%v: Cache '%s' recorded a cache hit on cache key '%v'.", event.CreationTime(), event.CacheName(), event.CacheKey())
		if event.String() != expectedMsg {
			t.Errorf("Expected event message '%s', but got '%s'", expectedMsg, event)
		}
		hits.Add(1)
	}
	onCacheMiss := func(event cache.MissEvent) {
		if event.EventType() != cache.OnMiss {
			t.Errorf("Expected event type CACHE_MISS, but got '%s'", event.EventType())
		}
		expectedMsg := fmt.Sprintf("%v: Cache '%s' recorded a cache miss on cache key '%v'.", event.CreationTime(), event.CacheName(), event.CacheKey())
		if event.String() != expectedMsg {
			t.Errorf("Expected event message '%s', but got '%s'", expectedMsg, event)
		}
		misses.Add(1)
	}
	eventListener.OnCacheHit(onCacheHit).OnCacheMiss(onCacheMiss)

	// fail with no error, max retries exceeded
	fn := func(key string) (string, error) {
		return key + randString(4), nil
	}
	decoratedFn := cache.DecorateFunction(ch, fn)

	var wg sync.WaitGroup
	var ret1, ret2 string
	wg.Add(1)
	go func() {
		ret1, _ = decoratedFn("notOK")
		wg.Done()
	}()
	wg.Add(1)
	go func() {
		ret2, _ = decoratedFn("notOK")
		wg.Done()
	}()
	wg.Wait()

	if ret1 != ret2 {
		t.Errorf("Expected return cached value, but got '%s' and '%s'", ret1, ret2)
	}

	time.Sleep(time.Second * 2)
	ret3, _ := decoratedFn("notOK")
	if ret1 == ret3 {
		t.Errorf("Expected return new value, but got '%s'", ret3)
	}

	if ch.Metrics().NumberOfCacheHits() != 1 {
		t.Errorf("Expected hits 1, but got '%d'", ch.Metrics().NumberOfCacheHits())
	}
	if ch.Metrics().NumberOfCacheMisses() != 2 {
		t.Errorf("Expected misses 2, but got '%d'", ch.Metrics().NumberOfCacheMisses())
	}

	time.Sleep(time.Second * 2)
	if hits.Load() != 1 {
		t.Errorf("Expected 1 hit call, but got '%d'", hits.Load())
	}
	if misses.Load() != 2 {
		t.Errorf("Expected 2 miss calls, but got '%d'", misses.Load())
	}
	eventListener.Dismiss(onCacheHit).Dismiss(onCacheMiss)
}
