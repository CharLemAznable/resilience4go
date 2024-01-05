package cache

import (
	"github.com/CharLemAznable/gogo/ext"
	"github.com/CharLemAznable/gogo/fn"
)

type EventListener interface {
	OnCacheHitFunc(func(HitEvent)) EventListener
	OnCacheMissFunc(func(MissEvent)) EventListener
	DismissCacheHitFunc(func(HitEvent)) EventListener
	DismissCacheMissFunc(func(MissEvent)) EventListener

	OnCacheHit(fn.Consumer[HitEvent]) EventListener
	OnCacheMiss(fn.Consumer[MissEvent]) EventListener
	DismissCacheHit(fn.Consumer[HitEvent]) EventListener
	DismissCacheMiss(fn.Consumer[MissEvent]) EventListener
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
	return listener.OnCacheHit(fn.ConsumerOf(consumer))
}

func (listener *eventListener) OnCacheMissFunc(consumer func(MissEvent)) EventListener {
	return listener.OnCacheMiss(fn.ConsumerOf(consumer))
}

func (listener *eventListener) DismissCacheHitFunc(consumer func(HitEvent)) EventListener {
	return listener.DismissCacheHit(fn.ConsumerOf(consumer))
}

func (listener *eventListener) DismissCacheMissFunc(consumer func(MissEvent)) EventListener {
	return listener.DismissCacheMiss(fn.ConsumerOf(consumer))
}

func (listener *eventListener) OnCacheHit(consumer fn.Consumer[HitEvent]) EventListener {
	listener.onCacheHit.AppendConsumer(consumer)
	return listener
}

func (listener *eventListener) OnCacheMiss(consumer fn.Consumer[MissEvent]) EventListener {
	listener.onCacheMiss.AppendConsumer(consumer)
	return listener
}

func (listener *eventListener) DismissCacheHit(consumer fn.Consumer[HitEvent]) EventListener {
	listener.onCacheHit.RemoveConsumer(consumer)
	return listener
}

func (listener *eventListener) DismissCacheMiss(consumer fn.Consumer[MissEvent]) EventListener {
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
