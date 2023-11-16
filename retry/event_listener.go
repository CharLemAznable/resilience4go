package retry

import (
	"github.com/CharLemAznable/resilience4go/utils"
	"sync"
)

type EventListener interface {
	OnSuccess(func(SuccessEvent)) EventListener
	OnRetry(func(RetryEvent)) EventListener
	OnError(func(ErrorEvent)) EventListener
	Dismiss(any) EventListener
	consumeEvent(Event)
}

func newEventListener() EventListener {
	return &eventListener{
		onSuccess: make([]func(SuccessEvent), 0),
		onRetry:   make([]func(RetryEvent), 0),
		onError:   make([]func(ErrorEvent), 0),
	}
}

type eventListener struct {
	sync.RWMutex
	onSuccess []func(SuccessEvent)
	onRetry   []func(RetryEvent)
	onError   []func(ErrorEvent)
}

func (listener *eventListener) OnSuccess(consumer func(SuccessEvent)) EventListener {
	listener.Lock()
	defer listener.Unlock()
	listener.onSuccess = utils.AppendElementUnique(listener.onSuccess, consumer)
	return listener
}

func (listener *eventListener) OnRetry(consumer func(RetryEvent)) EventListener {
	listener.Lock()
	defer listener.Unlock()
	listener.onRetry = utils.AppendElementUnique(listener.onRetry, consumer)
	return listener
}

func (listener *eventListener) OnError(consumer func(ErrorEvent)) EventListener {
	listener.Lock()
	defer listener.Unlock()
	listener.onError = utils.AppendElementUnique(listener.onError, consumer)
	return listener
}

func (listener *eventListener) Dismiss(consumer any) EventListener {
	listener.Lock()
	defer listener.Unlock()
	if c, ok := consumer.(func(SuccessEvent)); ok {
		listener.onSuccess = utils.RemoveElementByValue(listener.onSuccess, c)
	}
	if c, ok := consumer.(func(RetryEvent)); ok {
		listener.onRetry = utils.RemoveElementByValue(listener.onRetry, c)
	}
	if c, ok := consumer.(func(ErrorEvent)); ok {
		listener.onError = utils.RemoveElementByValue(listener.onError, c)
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
		case *retryEvent:
			utils.ConsumeEvent(listener.onRetry, RetryEvent(e))
		case *errorEvent:
			utils.ConsumeEvent(listener.onError, ErrorEvent(e))
		}
	}()
}
