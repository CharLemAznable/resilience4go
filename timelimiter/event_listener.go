package timelimiter

import (
	"github.com/CharLemAznable/resilience4go/utils"
	"sync"
)

type EventConsumer func(Event)

type EventListener interface {
	OnSuccess(EventConsumer) EventListener
	OnTimeout(EventConsumer) EventListener
	OnFailure(EventConsumer) EventListener
	Dismiss(EventConsumer) EventListener
	HasConsumer() bool
	consumeEvent(Event)
}

func newEventListener() EventListener {
	return &eventListener{
		onSuccess: make([]EventConsumer, 0),
		onTimeout: make([]EventConsumer, 0),
		onFailure: make([]EventConsumer, 0),
		slices:    utils.NewSlicesWithPointer[EventConsumer](),
	}
}

type eventListener struct {
	mutex     sync.RWMutex
	onSuccess []EventConsumer
	onTimeout []EventConsumer
	onFailure []EventConsumer
	slices    utils.Slices[EventConsumer]
}

func (listener *eventListener) OnSuccess(consumer EventConsumer) EventListener {
	listener.mutex.Lock()
	defer listener.mutex.Unlock()
	listener.onSuccess = listener.slices.AppendElementUnique(listener.onSuccess, consumer)
	return listener
}

func (listener *eventListener) OnTimeout(consumer EventConsumer) EventListener {
	listener.mutex.Lock()
	defer listener.mutex.Unlock()
	listener.onTimeout = listener.slices.AppendElementUnique(listener.onTimeout, consumer)
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
	listener.onTimeout = listener.slices.RemoveElementByValue(listener.onTimeout, consumer)
	listener.onFailure = listener.slices.RemoveElementByValue(listener.onFailure, consumer)
	return listener
}

func (listener *eventListener) HasConsumer() bool {
	listener.mutex.RLock()
	defer listener.mutex.RUnlock()
	return len(listener.onSuccess) > 0 || len(listener.onTimeout) > 0 || len(listener.onFailure) > 0
}

func (listener *eventListener) consumeEvent(event Event) {
	if !listener.HasConsumer() {
		return
	}
	listener.mutex.RLock()
	defer listener.mutex.RUnlock()
	var consumers []EventConsumer
	switch event.EventType() {
	case SUCCESS:
		consumers = listener.onSuccess
	case TIMEOUT:
		consumers = listener.onTimeout
	case FAILURE:
		consumers = listener.onFailure
	}
	for _, consumer := range consumers {
		go consumer(event)
	}
}
