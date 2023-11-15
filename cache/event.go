package cache

import (
	"fmt"
	"time"
)

type EventType string

const (
	OnHit  EventType = "CACHE_HIT"
	OnMiss EventType = "CACHE_MISS"
)

type Event interface {
	fmt.Stringer
	CacheName() string
	CreationTime() time.Time
	CacheKey() any
	EventType() EventType
}

func newCacheHitEvent(cacheName string, key any) Event {
	return &hitEvent{event{cacheName: cacheName, creationTime: time.Now(), cacheKey: key}}
}

func newCacheMissEvent(cacheName string, key any) Event {
	return &missEvent{event{cacheName: cacheName, creationTime: time.Now(), cacheKey: key}}
}

type event struct {
	cacheName    string
	creationTime time.Time
	cacheKey     any
}

func (e *event) CacheName() string {
	return e.cacheName
}

func (e *event) CreationTime() time.Time {
	return e.creationTime
}

func (e *event) CacheKey() any {
	return e.cacheKey
}

type hitEvent struct {
	event
}

func (e *hitEvent) EventType() EventType {
	return OnHit
}

func (e *hitEvent) String() string {
	return fmt.Sprintf(
		"%v: Cache '%s' recorded a cache hit on cache key '%v'.",
		e.creationTime, e.cacheName, e.cacheKey)
}

type missEvent struct {
	event
}

func (e *missEvent) EventType() EventType {
	return OnMiss
}

func (e *missEvent) String() string {
	return fmt.Sprintf(
		"%v: Cache '%s' recorded a cache miss on cache key '%v'.",
		e.creationTime, e.cacheName, e.cacheKey)
}
