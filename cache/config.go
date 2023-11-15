package cache

import (
	"fmt"
	"github.com/cespare/xxhash/v2"
	"github.com/dgraph-io/ristretto/z"
	"time"
)

type Config struct {
	capacity  int64
	itemTTL   time.Duration
	keyToHash func(key any) (keyHash uint64, conflictHash uint64)
}

func (config *Config) String() string {
	return fmt.Sprintf(
		"CacheConfig"+
			" {capacity=%d, itemTTL=%v"+
			", keyToHash %T[%v]}",
		config.capacity, config.itemTTL,
		config.keyToHash, any(config.keyToHash))
}

func (config *Config) keyToHashFn(key any) (uint64, uint64) {
	if config.keyToHash != nil {
		return config.keyToHash(key)
	}
	return DefaultKeyToHash(key)
}

type ConfigBuilder func(*Config)

func WithCapacity(capacity int64) ConfigBuilder {
	return func(config *Config) {
		config.capacity = capacity
	}
}

func WithItemTTL(itemTTL time.Duration) ConfigBuilder {
	return func(config *Config) {
		config.itemTTL = itemTTL
	}
}

func WithKeyToHash(function func(any) (uint64, uint64)) ConfigBuilder {
	return func(config *Config) {
		config.keyToHash = function
	}
}

const DefaultCapacity int64 = 10000
const DefaultItemTTL = time.Minute * 5

func DefaultKeyToHash(key any) (uint64, uint64) {
	fmtKey := fmt.Sprintf("%v", key)
	return z.MemHashString(fmtKey), xxhash.Sum64String(fmtKey)
}

func defaultConfig() *Config {
	return &Config{
		capacity:  DefaultCapacity,
		itemTTL:   DefaultItemTTL,
		keyToHash: DefaultKeyToHash,
	}
}
