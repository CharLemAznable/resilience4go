package cache

import (
	"github.com/CharLemAznable/gogo/lang"
	"github.com/dgraph-io/ristretto"
	"sync"
)

type Cache[K any, V any] interface {
	Name() string
	Metrics() Metrics
	EventListener() EventListener
	WithMarshalFn(func(V) any, func(any) V) Cache[K, V]
	GetOrLoad(key K, loader func(K) (V, error)) (V, error)
}

func NewCache[K any, V any](name string, configs ...ConfigBuilder) Cache[K, V] {
	config := defaultConfig()
	for _, cfg := range configs {
		cfg(config)
	}
	ristrettoCache, err := ristretto.NewCache(&ristretto.Config{
		NumCounters:        config.capacity * 10,
		MaxCost:            config.capacity,
		BufferItems:        64,
		Metrics:            true,
		KeyToHash:          config.keyToHashFn,
		IgnoreInternalCost: true,
	})
	lang.PanicIfError(err)
	c := &cache[K, V]{
		name:           name,
		config:         config,
		locks:          make([]*sync.Mutex, int(numShards)),
		ristrettoCache: ristrettoCache,
		metrics:        newMetrics(ristrettoCache.Metrics),
		eventListener:  newEventListener(),
	}
	for i := range c.locks {
		c.locks[i] = &sync.Mutex{}
	}
	return c
}

const numShards uint64 = 256

type cache[K any, V any] struct {
	name           string
	config         *Config
	locks          []*sync.Mutex
	ristrettoCache *ristretto.Cache
	metrics        *metrics
	eventListener  *eventListener

	marshalFn   func(V) any
	unmarshalFn func(any) V
}

func (c *cache[K, V]) Name() string {
	return c.name
}

func (c *cache[K, V]) Metrics() Metrics {
	return c.metrics
}

func (c *cache[K, V]) EventListener() EventListener {
	return c.eventListener
}

func (c *cache[K, V]) WithMarshalFn(marshalFn func(V) any, unmarshalFn func(any) V) Cache[K, V] {
	c.marshalFn, c.unmarshalFn = marshalFn, unmarshalFn
	return c
}

func (c *cache[K, V]) GetOrLoad(key K, loader func(K) (V, error)) (V, error) {
	keyHash, _ := c.config.keyToHashFn(key)
	lockIdx := keyHash % numShards
	c.locks[lockIdx].Lock()
	defer c.locks[lockIdx].Unlock()

	if v, found := c.ristrettoCache.Get(key); found {
		c.eventListener.consumeEvent(newCacheHitEvent(c.name, key))
		vv, err := lang.Cast[*valueWithError](v)
		lang.PanicIfError(err)
		vvv := vv.value
		if c.unmarshalFn != nil {
			vvv = c.unmarshalFn(vv.value)
		}
		value, err := lang.CastOrZero[V](vvv)
		lang.PanicIfError(err)
		return value, vv.error
	}
	c.eventListener.consumeEvent(newCacheMissEvent(c.name, key))
	value, err := loader(key)
	if c.config.cacheResultPredicateFn(value, err) {
		var vvv any = value
		if c.marshalFn != nil {
			vvv = c.marshalFn(value)
		}
		vv := &valueWithError{value: vvv, error: err}
		c.ristrettoCache.SetWithTTL(key, vv, 1, c.config.itemTTL)
		c.ristrettoCache.Wait()
	}
	return value, err
}

type valueWithError struct {
	value any
	error error
}
