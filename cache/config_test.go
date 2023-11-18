package cache_test

import (
	"fmt"
	"github.com/CharLemAznable/resilience4go/cache"
	"testing"
	"time"
)

func TestConfig_String(t *testing.T) {
	config := &cache.Config{}
	cache.WithCapacity(10)(config)
	cache.WithItemTTL(time.Second * 5)(config)
	keyToHash := func(key any) (uint64, uint64) {
		return 0, 0
	}
	cache.WithKeyToHash(keyToHash)(config)
	cacheResultPredicate := func(_ any, _ error) bool {
		return false
	}
	cache.WithCacheResultPredicate(cacheResultPredicate)(config)
	expected := fmt.Sprintf("CacheConfig {capacity=10, itemTTL=5s, keyToHash %T[%v], cacheResultPredicate %T[%v]}",
		keyToHash, any(keyToHash), cacheResultPredicate, any(cacheResultPredicate))
	result := fmt.Sprintf("%v", config)
	if result != expected {
		t.Errorf("Expected %s, but got %s", expected, result)
	}
}
