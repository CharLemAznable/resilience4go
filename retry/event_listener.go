package retry

import (
	"github.com/CharLemAznable/ge"
	"sync"
)

type EventListener interface {
	OnSuccess(func(SuccessEvent)) EventListener
	OnRetry(func(RetryEvent)) EventListener
	OnError(func(ErrorEvent)) EventListener
	Dismiss(any) EventListener
}

func newEventListener() *eventListener {
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
	listener.onSuccess = ge.AppendElementUnique(listener.onSuccess, consumer)
	return listener
}

func (listener *eventListener) OnRetry(consumer func(RetryEvent)) EventListener {
	listener.Lock()
	defer listener.Unlock()
	listener.onRetry = ge.AppendElementUnique(listener.onRetry, consumer)
	return listener
}

func (listener *eventListener) OnError(consumer func(ErrorEvent)) EventListener {
	listener.Lock()
	defer listener.Unlock()
	listener.onError = ge.AppendElementUnique(listener.onError, consumer)
	return listener
}

func (listener *eventListener) Dismiss(consumer any) EventListener {
	listener.Lock()
	defer listener.Unlock()
	switch c := consumer.(type) {
	case func(SuccessEvent):
		listener.onSuccess = ge.RemoveElementByValue(listener.onSuccess, c)
	case func(RetryEvent):
		listener.onRetry = ge.RemoveElementByValue(listener.onRetry, c)
	case func(ErrorEvent):
		listener.onError = ge.RemoveElementByValue(listener.onError, c)
	}
	return listener
}

func (listener *eventListener) consumeEvent(event Event) {
	go func() {
		listener.RLock()
		defer listener.RUnlock()
		switch e := event.(type) {
		case *successEvent:
			ge.ConsumeEach(listener.onSuccess, SuccessEvent(e))
		case *retryEvent:
			ge.ConsumeEach(listener.onRetry, RetryEvent(e))
		case *errorEvent:
			ge.ConsumeEach(listener.onError, ErrorEvent(e))
		}
	}()
}
