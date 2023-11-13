package circuitbreaker

import (
	"github.com/CharLemAznable/resilience4go/utils"
	"sync"
)

type EventConsumer func(Event)

type EventListener interface {
	OnSuccess(EventConsumer) EventListener
	OnError(EventConsumer) EventListener
	OnNotPermitted(EventConsumer) EventListener
	OnStateTransition(EventConsumer) EventListener
	OnFailureRateExceeded(EventConsumer) EventListener
	OnSlowCallRateExceeded(EventConsumer) EventListener
	Dismiss(EventConsumer) EventListener
	HasConsumer() bool
	consumeEvent(Event)
}

func newEventListener() EventListener {
	return &eventListener{
		onSuccess:              make([]EventConsumer, 0),
		onError:                make([]EventConsumer, 0),
		onNotPermitted:         make([]EventConsumer, 0),
		onStateTransition:      make([]EventConsumer, 0),
		onFailureRateExceeded:  make([]EventConsumer, 0),
		onSlowCallRateExceeded: make([]EventConsumer, 0),
		slices:                 utils.NewSlicesWithPointer[EventConsumer](),
	}
}

type eventListener struct {
	mutex                  sync.RWMutex
	onSuccess              []EventConsumer
	onError                []EventConsumer
	onNotPermitted         []EventConsumer
	onStateTransition      []EventConsumer
	onFailureRateExceeded  []EventConsumer
	onSlowCallRateExceeded []EventConsumer
	slices                 utils.Slices[EventConsumer]
}

func (listener *eventListener) OnSuccess(consumer EventConsumer) EventListener {
	listener.mutex.Lock()
	defer listener.mutex.Unlock()
	listener.onSuccess = listener.slices.AppendElementUnique(listener.onSuccess, consumer)
	return listener
}

func (listener *eventListener) OnError(consumer EventConsumer) EventListener {
	listener.mutex.Lock()
	defer listener.mutex.Unlock()
	listener.onError = listener.slices.AppendElementUnique(listener.onError, consumer)
	return listener
}

func (listener *eventListener) OnNotPermitted(consumer EventConsumer) EventListener {
	listener.mutex.Lock()
	defer listener.mutex.Unlock()
	listener.onNotPermitted = listener.slices.AppendElementUnique(listener.onNotPermitted, consumer)
	return listener
}

func (listener *eventListener) OnStateTransition(consumer EventConsumer) EventListener {
	listener.mutex.Lock()
	defer listener.mutex.Unlock()
	listener.onStateTransition = listener.slices.AppendElementUnique(listener.onStateTransition, consumer)
	return listener
}

func (listener *eventListener) OnFailureRateExceeded(consumer EventConsumer) EventListener {
	listener.mutex.Lock()
	defer listener.mutex.Unlock()
	listener.onFailureRateExceeded = listener.slices.AppendElementUnique(listener.onFailureRateExceeded, consumer)
	return listener
}

func (listener *eventListener) OnSlowCallRateExceeded(consumer EventConsumer) EventListener {
	listener.mutex.Lock()
	defer listener.mutex.Unlock()
	listener.onSlowCallRateExceeded = listener.slices.AppendElementUnique(listener.onSlowCallRateExceeded, consumer)
	return listener
}

func (listener *eventListener) Dismiss(consumer EventConsumer) EventListener {
	listener.mutex.Lock()
	defer listener.mutex.Unlock()
	listener.onSuccess = listener.slices.RemoveElementByValue(listener.onSuccess, consumer)
	listener.onError = listener.slices.RemoveElementByValue(listener.onError, consumer)
	listener.onNotPermitted = listener.slices.RemoveElementByValue(listener.onNotPermitted, consumer)
	listener.onStateTransition = listener.slices.RemoveElementByValue(listener.onStateTransition, consumer)
	listener.onFailureRateExceeded = listener.slices.RemoveElementByValue(listener.onFailureRateExceeded, consumer)
	listener.onSlowCallRateExceeded = listener.slices.RemoveElementByValue(listener.onSlowCallRateExceeded, consumer)
	return listener
}

func (listener *eventListener) HasConsumer() bool {
	listener.mutex.RLock()
	defer listener.mutex.RUnlock()
	return len(listener.onSuccess) > 0 || len(listener.onError) > 0 || len(listener.onNotPermitted) > 0 ||
		len(listener.onStateTransition) > 0 || len(listener.onFailureRateExceeded) > 0 || len(listener.onSlowCallRateExceeded) > 0
}

func (listener *eventListener) consumeEvent(event Event) {
	if !listener.HasConsumer() {
		return
	}
	listener.mutex.RLock()
	defer listener.mutex.RUnlock()
	var consumers []EventConsumer
	switch event.EventType() {
	case Success:
		consumers = listener.onSuccess
	case Error:
		consumers = listener.onError
	case NotPermitted:
		consumers = listener.onNotPermitted
	case StateTransition:
		consumers = listener.onStateTransition
	case FailureRateExceeded:
		consumers = listener.onFailureRateExceeded
	case SlowCallRateExceeded:
		consumers = listener.onSlowCallRateExceeded
	}
	for _, consumer := range consumers {
		go consumer(event)
	}
}
