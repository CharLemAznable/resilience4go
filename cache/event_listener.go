package cache

import (
	"github.com/CharLemAznable/ge"
	"sync"
)

type EventListener interface {
	OnCacheHitFunc(func(HitEvent)) EventListener
	OnCacheMissFunc(func(MissEvent)) EventListener
	DismissCacheHitFunc(func(HitEvent)) EventListener
	DismissCacheMissFunc(func(MissEvent)) EventListener

	OnCacheHit(ge.Action[HitEvent]) EventListener
	OnCacheMiss(ge.Action[MissEvent]) EventListener
	DismissCacheHit(ge.Action[HitEvent]) EventListener
	DismissCacheMiss(ge.Action[MissEvent]) EventListener
}

func newEventListener() *eventListener {
	return &eventListener{
		onCacheHit:  make([]ge.Action[HitEvent], 0),
		onCacheMiss: make([]ge.Action[MissEvent], 0),
	}
}

type eventListener struct {
	sync.RWMutex
	onCacheHit  []ge.Action[HitEvent]
	onCacheMiss []ge.Action[MissEvent]
}

func (listener *eventListener) OnCacheHitFunc(consumer func(HitEvent)) EventListener {
	return listener.OnCacheHit(ge.ActionFunc[HitEvent](consumer))
}

func (listener *eventListener) OnCacheMissFunc(consumer func(MissEvent)) EventListener {
	return listener.OnCacheMiss(ge.ActionFunc[MissEvent](consumer))
}

func (listener *eventListener) DismissCacheHitFunc(consumer func(HitEvent)) EventListener {
	return listener.DismissCacheHit(ge.ActionFunc[HitEvent](consumer))
}

func (listener *eventListener) DismissCacheMissFunc(consumer func(MissEvent)) EventListener {
	return listener.DismissCacheMiss(ge.ActionFunc[MissEvent](consumer))
}

func (listener *eventListener) OnCacheHit(consumer ge.Action[HitEvent]) EventListener {
	listener.Lock()
	defer listener.Unlock()
	listener.onCacheHit = ge.AppendElementUnique(listener.onCacheHit, consumer)
	return listener
}

func (listener *eventListener) OnCacheMiss(consumer ge.Action[MissEvent]) EventListener {
	listener.Lock()
	defer listener.Unlock()
	listener.onCacheMiss = ge.AppendElementUnique(listener.onCacheMiss, consumer)
	return listener
}

func (listener *eventListener) DismissCacheHit(consumer ge.Action[HitEvent]) EventListener {
	listener.Lock()
	defer listener.Unlock()
	listener.onCacheHit = ge.RemoveElementByValue(listener.onCacheHit, consumer)
	return listener
}

func (listener *eventListener) DismissCacheMiss(consumer ge.Action[MissEvent]) EventListener {
	listener.Lock()
	defer listener.Unlock()
	listener.onCacheMiss = ge.RemoveElementByValue(listener.onCacheMiss, consumer)
	return listener
}

func (listener *eventListener) consumeEvent(event Event) {
	go func() {
		listener.RLock()
		defer listener.RUnlock()
		switch e := event.(type) {
		case *hitEvent:
			ge.ForEach(listener.onCacheHit, HitEvent(e))
		case *missEvent:
			ge.ForEach(listener.onCacheMiss, MissEvent(e))
		}
	}()
}
