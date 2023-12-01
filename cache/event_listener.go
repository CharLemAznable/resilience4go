package cache

import (
	"github.com/CharLemAznable/resilience4go/utils"
	"sync"
)

type EventListener interface {
	OnCacheHit(func(HitEvent)) EventListener
	OnCacheMiss(func(MissEvent)) EventListener
	Dismiss(any) EventListener
}

func newEventListener() *eventListener {
	return &eventListener{
		onCacheHit:  make([]func(HitEvent), 0),
		onCacheMiss: make([]func(MissEvent), 0),
	}
}

type eventListener struct {
	sync.RWMutex
	onCacheHit  []func(HitEvent)
	onCacheMiss []func(MissEvent)
}

func (listener *eventListener) OnCacheHit(consumer func(HitEvent)) EventListener {
	listener.Lock()
	defer listener.Unlock()
	listener.onCacheHit = utils.AppendElementUnique(listener.onCacheHit, consumer)
	return listener
}

func (listener *eventListener) OnCacheMiss(consumer func(MissEvent)) EventListener {
	listener.Lock()
	defer listener.Unlock()
	listener.onCacheMiss = utils.AppendElementUnique(listener.onCacheMiss, consumer)
	return listener
}

func (listener *eventListener) Dismiss(consumer any) EventListener {
	listener.Lock()
	defer listener.Unlock()
	switch c := consumer.(type) {
	case func(HitEvent):
		listener.onCacheHit = utils.RemoveElementByValue(listener.onCacheHit, c)
	case func(MissEvent):
		listener.onCacheMiss = utils.RemoveElementByValue(listener.onCacheMiss, c)
	}
	return listener
}

func (listener *eventListener) consumeEvent(event Event) {
	go func() {
		listener.RLock()
		defer listener.RUnlock()
		switch e := event.(type) {
		case *hitEvent:
			utils.ConsumeEvent(listener.onCacheHit, HitEvent(e))
		case *missEvent:
			utils.ConsumeEvent(listener.onCacheMiss, MissEvent(e))
		}
	}()
}
