package retry

import (
	"github.com/CharLemAznable/ge"
	"sync"
)

type EventListener interface {
	OnSuccessFunc(func(SuccessEvent)) EventListener
	OnRetryFunc(func(RetryEvent)) EventListener
	OnErrorFunc(func(ErrorEvent)) EventListener
	DismissSuccessFunc(func(SuccessEvent)) EventListener
	DismissRetryFunc(func(RetryEvent)) EventListener
	DismissErrorFunc(func(ErrorEvent)) EventListener

	OnSuccess(ge.Action[SuccessEvent]) EventListener
	OnRetry(ge.Action[RetryEvent]) EventListener
	OnError(ge.Action[ErrorEvent]) EventListener
	DismissSuccess(ge.Action[SuccessEvent]) EventListener
	DismissRetry(ge.Action[RetryEvent]) EventListener
	DismissError(ge.Action[ErrorEvent]) EventListener
}

func newEventListener() *eventListener {
	return &eventListener{
		onSuccess: make([]ge.Action[SuccessEvent], 0),
		onRetry:   make([]ge.Action[RetryEvent], 0),
		onError:   make([]ge.Action[ErrorEvent], 0),
	}
}

type eventListener struct {
	sync.RWMutex
	onSuccess []ge.Action[SuccessEvent]
	onRetry   []ge.Action[RetryEvent]
	onError   []ge.Action[ErrorEvent]
}

func (listener *eventListener) OnSuccessFunc(consumer func(SuccessEvent)) EventListener {
	return listener.OnSuccess(ge.ActionFunc[SuccessEvent](consumer))
}

func (listener *eventListener) OnRetryFunc(consumer func(RetryEvent)) EventListener {
	return listener.OnRetry(ge.ActionFunc[RetryEvent](consumer))
}

func (listener *eventListener) OnErrorFunc(consumer func(ErrorEvent)) EventListener {
	return listener.OnError(ge.ActionFunc[ErrorEvent](consumer))
}

func (listener *eventListener) DismissSuccessFunc(consumer func(SuccessEvent)) EventListener {
	return listener.DismissSuccess(ge.ActionFunc[SuccessEvent](consumer))
}

func (listener *eventListener) DismissRetryFunc(consumer func(RetryEvent)) EventListener {
	return listener.DismissRetry(ge.ActionFunc[RetryEvent](consumer))
}

func (listener *eventListener) DismissErrorFunc(consumer func(ErrorEvent)) EventListener {
	return listener.DismissError(ge.ActionFunc[ErrorEvent](consumer))
}

func (listener *eventListener) OnSuccess(action ge.Action[SuccessEvent]) EventListener {
	listener.Lock()
	defer listener.Unlock()
	listener.onSuccess = ge.AppendElementUnique(listener.onSuccess, action)
	return listener
}

func (listener *eventListener) OnRetry(action ge.Action[RetryEvent]) EventListener {
	listener.Lock()
	defer listener.Unlock()
	listener.onRetry = ge.AppendElementUnique(listener.onRetry, action)
	return listener
}

func (listener *eventListener) OnError(action ge.Action[ErrorEvent]) EventListener {
	listener.Lock()
	defer listener.Unlock()
	listener.onError = ge.AppendElementUnique(listener.onError, action)
	return listener
}

func (listener *eventListener) DismissSuccess(action ge.Action[SuccessEvent]) EventListener {
	listener.Lock()
	defer listener.Unlock()
	listener.onSuccess = ge.RemoveElementByValue(listener.onSuccess, action)
	return listener
}

func (listener *eventListener) DismissRetry(action ge.Action[RetryEvent]) EventListener {
	listener.Lock()
	defer listener.Unlock()
	listener.onRetry = ge.RemoveElementByValue(listener.onRetry, action)
	return listener
}

func (listener *eventListener) DismissError(action ge.Action[ErrorEvent]) EventListener {
	listener.Lock()
	defer listener.Unlock()
	listener.onError = ge.RemoveElementByValue(listener.onError, action)
	return listener
}

func (listener *eventListener) consumeEvent(event Event) {
	go func() {
		listener.RLock()
		defer listener.RUnlock()
		switch e := event.(type) {
		case *successEvent:
			ge.ForEach(listener.onSuccess, SuccessEvent(e))
		case *retryEvent:
			ge.ForEach(listener.onRetry, RetryEvent(e))
		case *errorEvent:
			ge.ForEach(listener.onError, ErrorEvent(e))
		}
	}()
}
