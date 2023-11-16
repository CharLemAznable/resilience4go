package timelimiter

import (
	"github.com/CharLemAznable/resilience4go/utils"
	"sync"
)

type EventListener interface {
	OnSuccess(func(SuccessEvent)) EventListener
	OnTimeout(func(TimeoutEvent)) EventListener
	OnFailure(func(FailureEvent)) EventListener
	Dismiss(any) EventListener
	consumeEvent(Event)
}

func newEventListener() EventListener {
	return &eventListener{
		onSuccess: make([]func(SuccessEvent), 0),
		onTimeout: make([]func(TimeoutEvent), 0),
		onFailure: make([]func(FailureEvent), 0),
	}
}

type eventListener struct {
	sync.RWMutex
	onSuccess []func(SuccessEvent)
	onTimeout []func(TimeoutEvent)
	onFailure []func(FailureEvent)
}

func (listener *eventListener) OnSuccess(consumer func(SuccessEvent)) EventListener {
	listener.Lock()
	defer listener.Unlock()
	listener.onSuccess = utils.AppendElementUnique(listener.onSuccess, consumer)
	return listener
}

func (listener *eventListener) OnTimeout(consumer func(TimeoutEvent)) EventListener {
	listener.Lock()
	defer listener.Unlock()
	listener.onTimeout = utils.AppendElementUnique(listener.onTimeout, consumer)
	return listener
}

func (listener *eventListener) OnFailure(consumer func(FailureEvent)) EventListener {
	listener.Lock()
	defer listener.Unlock()
	listener.onFailure = utils.AppendElementUnique(listener.onFailure, consumer)
	return listener
}

func (listener *eventListener) Dismiss(consumer any) EventListener {
	listener.Lock()
	defer listener.Unlock()
	if c, ok := consumer.(func(SuccessEvent)); ok {
		listener.onSuccess = utils.RemoveElementByValue(listener.onSuccess, c)
	}
	if c, ok := consumer.(func(TimeoutEvent)); ok {
		listener.onTimeout = utils.RemoveElementByValue(listener.onTimeout, c)
	}
	if c, ok := consumer.(func(FailureEvent)); ok {
		listener.onFailure = utils.RemoveElementByValue(listener.onFailure, c)
	}
	return listener
}

func (listener *eventListener) consumeEvent(event Event) {
	go func() {
		listener.RLock()
		defer listener.RUnlock()
		switch e := event.(type) {
		case *successEvent:
			utils.ConsumeEvent(listener.onSuccess, SuccessEvent(e))
		case *timeoutEvent:
			utils.ConsumeEvent(listener.onTimeout, TimeoutEvent(e))
		case *failureEvent:
			utils.ConsumeEvent(listener.onFailure, FailureEvent(e))
		}
	}()
}
