package cache

import (
	"github.com/CharLemAznable/gogo/ext"
	. "github.com/CharLemAznable/gogo/fn"
)

type EventListener interface {
	OnCacheHitFunc(func(HitEvent)) EventListener
	OnCacheMissFunc(func(MissEvent)) EventListener
	DismissCacheHitFunc(func(HitEvent)) EventListener
	DismissCacheMissFunc(func(MissEvent)) EventListener

	OnCacheHit(Consumer[HitEvent]) EventListener
	OnCacheMiss(Consumer[MissEvent]) EventListener
	DismissCacheHit(Consumer[HitEvent]) EventListener
	DismissCacheMiss(Consumer[MissEvent]) EventListener
}

func newEventListener() *eventListener {
	return &eventListener{
		onCacheHit:  ext.NewConsumers[HitEvent](),
		onCacheMiss: ext.NewConsumers[MissEvent](),
	}
}

type eventListener struct {
	onCacheHit  ext.Consumers[HitEvent]
	onCacheMiss ext.Consumers[MissEvent]
}

func (listener *eventListener) OnCacheHitFunc(consumer func(HitEvent)) EventListener {
	return listener.OnCacheHit(ConsumerOf(consumer))
}

func (listener *eventListener) OnCacheMissFunc(consumer func(MissEvent)) EventListener {
	return listener.OnCacheMiss(ConsumerOf(consumer))
}

func (listener *eventListener) DismissCacheHitFunc(consumer func(HitEvent)) EventListener {
	return listener.DismissCacheHit(ConsumerOf(consumer))
}

func (listener *eventListener) DismissCacheMissFunc(consumer func(MissEvent)) EventListener {
	return listener.DismissCacheMiss(ConsumerOf(consumer))
}

func (listener *eventListener) OnCacheHit(consumer Consumer[HitEvent]) EventListener {
	listener.onCacheHit.AppendConsumer(consumer)
	return listener
}

func (listener *eventListener) OnCacheMiss(consumer Consumer[MissEvent]) EventListener {
	listener.onCacheMiss.AppendConsumer(consumer)
	return listener
}

func (listener *eventListener) DismissCacheHit(consumer Consumer[HitEvent]) EventListener {
	listener.onCacheHit.RemoveConsumer(consumer)
	return listener
}

func (listener *eventListener) DismissCacheMiss(consumer Consumer[MissEvent]) EventListener {
	listener.onCacheMiss.RemoveConsumer(consumer)
	return listener
}

func (listener *eventListener) consumeEvent(event Event) {
	go func() {
		switch e := event.(type) {
		case *hitEvent:
			listener.onCacheHit.Accept(HitEvent(e))
		case *missEvent:
			listener.onCacheMiss.Accept(MissEvent(e))
		}
	}()
}
