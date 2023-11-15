package cache

import (
	"github.com/CharLemAznable/resilience4go/utils"
	"sync"
)

type EventConsumer func(Event)

type EventListener interface {
	OnCacheHit(EventConsumer) EventListener
	OnCacheMiss(EventConsumer) EventListener
	Dismiss(EventConsumer) EventListener
	HasConsumer() bool
	consumeEvent(Event)
}

func newEventListener() EventListener {
	return &eventListener{
		onCacheHit:  make([]EventConsumer, 0),
		onCacheMiss: make([]EventConsumer, 0),
		slices:      utils.NewSlicesWithPointer[EventConsumer](),
	}
}

type eventListener struct {
	mutex       sync.RWMutex
	onCacheHit  []EventConsumer
	onCacheMiss []EventConsumer
	slices      utils.Slices[EventConsumer]
}

func (listener *eventListener) OnCacheHit(consumer EventConsumer) EventListener {
	listener.mutex.Lock()
	defer listener.mutex.Unlock()
	listener.onCacheHit = listener.slices.AppendElementUnique(listener.onCacheHit, consumer)
	return listener
}

func (listener *eventListener) OnCacheMiss(consumer EventConsumer) EventListener {
	listener.mutex.Lock()
	defer listener.mutex.Unlock()
	listener.onCacheMiss = listener.slices.AppendElementUnique(listener.onCacheMiss, consumer)
	return listener
}

func (listener *eventListener) Dismiss(consumer EventConsumer) EventListener {
	listener.mutex.Lock()
	defer listener.mutex.Unlock()
	listener.onCacheHit = listener.slices.RemoveElementByValue(listener.onCacheHit, consumer)
	listener.onCacheMiss = listener.slices.RemoveElementByValue(listener.onCacheMiss, consumer)
	return listener
}

func (listener *eventListener) HasConsumer() bool {
	listener.mutex.RLock()
	defer listener.mutex.RUnlock()
	return len(listener.onCacheHit) > 0 || len(listener.onCacheMiss) > 0
}

func (listener *eventListener) consumeEvent(event Event) {
	go func() {
		if !listener.HasConsumer() {
			return
		}
		listener.mutex.RLock()
		defer listener.mutex.RUnlock()
		var consumers []EventConsumer
		switch event.EventType() {
		case OnHit:
			consumers = listener.onCacheHit
		case OnMiss:
			consumers = listener.onCacheMiss
		}
		for _, consumer := range consumers {
			go consumer(event)
		}
	}()
}
