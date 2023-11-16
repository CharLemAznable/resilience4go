package cache

import (
	"github.com/CharLemAznable/resilience4go/utils"
	"sync"
)

type HitEventConsumer func(HitEvent)
type MissEventConsumer func(MissEvent)

type EventListener interface {
	OnCacheHit(HitEventConsumer) EventListener
	OnCacheMiss(MissEventConsumer) EventListener
	Dismiss(any) EventListener
	consumeEvent(Event)
}

func newEventListener() EventListener {
	return &eventListener{
		onCacheHit:        make([]HitEventConsumer, 0),
		onCacheHitSlices:  utils.NewSlicesWithPointer[HitEventConsumer](),
		onCacheMiss:       make([]MissEventConsumer, 0),
		onCacheMissSlices: utils.NewSlicesWithPointer[MissEventConsumer](),
	}
}

type eventListener struct {
	sync.RWMutex
	onCacheHit        []HitEventConsumer
	onCacheHitSlices  utils.Slices[HitEventConsumer]
	onCacheMiss       []MissEventConsumer
	onCacheMissSlices utils.Slices[MissEventConsumer]
}

func (listener *eventListener) OnCacheHit(consumer HitEventConsumer) EventListener {
	listener.Lock()
	defer listener.Unlock()
	listener.onCacheHit = listener.onCacheHitSlices.AppendElementUnique(listener.onCacheHit, consumer)
	return listener
}

func (listener *eventListener) OnCacheMiss(consumer MissEventConsumer) EventListener {
	listener.Lock()
	defer listener.Unlock()
	listener.onCacheMiss = listener.onCacheMissSlices.AppendElementUnique(listener.onCacheMiss, consumer)
	return listener
}

func (listener *eventListener) Dismiss(consumer any) EventListener {
	listener.Lock()
	defer listener.Unlock()
	if c, ok := consumer.(func(HitEvent)); ok {
		listener.onCacheHit = listener.onCacheHitSlices.RemoveElementByValue(listener.onCacheHit, c)
	}
	if c, ok := consumer.(func(MissEvent)); ok {
		listener.onCacheMiss = listener.onCacheMissSlices.RemoveElementByValue(listener.onCacheMiss, c)
	}
	return listener
}

func (listener *eventListener) consumeEvent(event Event) {
	go func() {
		listener.RLock()
		defer listener.RUnlock()
		switch e := event.(type) {
		case *hitEvent:
			for _, consumer := range listener.onCacheHit {
				go consumer(e)
			}
		case *missEvent:
			for _, consumer := range listener.onCacheMiss {
				go consumer(e)
			}
		}
	}()
}
