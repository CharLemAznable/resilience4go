package ratelimiter

import (
	"github.com/CharLemAznable/resilience4go/utils"
	"sync"
)

type EventConsumer func(Event)

type EventListener interface {
	OnSuccess(EventConsumer) EventListener
	OnFailure(EventConsumer) EventListener
	Dismiss(EventConsumer) EventListener
	HasConsumer() bool
	consumeEvent(Event)
}

func newEventListener() EventListener {
	return &eventListener{
		onSuccess: make([]EventConsumer, 0),
		onFailure: make([]EventConsumer, 0),
		slices:    utils.NewSlicesWithPointer[EventConsumer](),
	}
}

type eventListener struct {
	mutex     sync.RWMutex
	onSuccess []EventConsumer
	onFailure []EventConsumer
	slices    utils.Slices[EventConsumer]
}

func (listener *eventListener) OnSuccess(consumer EventConsumer) EventListener {
	listener.mutex.Lock()
	defer listener.mutex.Unlock()
	listener.onSuccess = listener.slices.AppendElementUnique(listener.onSuccess, consumer)
	return listener
}

func (listener *eventListener) OnFailure(consumer EventConsumer) EventListener {
	listener.mutex.Lock()
	defer listener.mutex.Unlock()
	listener.onFailure = listener.slices.AppendElementUnique(listener.onFailure, consumer)
	return listener
}

func (listener *eventListener) Dismiss(consumer EventConsumer) EventListener {
	listener.mutex.Lock()
	defer listener.mutex.Unlock()
	listener.onSuccess = listener.slices.RemoveElementByValue(listener.onSuccess, consumer)
	listener.onFailure = listener.slices.RemoveElementByValue(listener.onFailure, consumer)
	return listener
}

func (listener *eventListener) HasConsumer() bool {
	listener.mutex.RLock()
	defer listener.mutex.RUnlock()
	return len(listener.onSuccess) > 0 || len(listener.onFailure) > 0
}

func (listener *eventListener) consumeEvent(event Event) {
	if !listener.HasConsumer() {
		return
	}
	listener.mutex.RLock()
	defer listener.mutex.RUnlock()
	var consumers []EventConsumer
	switch event.EventType() {
	case SUCCESSFUL:
		consumers = listener.onSuccess
	case FAILED:
		consumers = listener.onFailure
	}
	for _, consumer := range consumers {
		go consumer(event)
	}
}
