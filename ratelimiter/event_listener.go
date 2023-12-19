package ratelimiter

import (
	"github.com/CharLemAznable/ge"
	"sync"
)

type EventListener interface {
	OnSuccessFunc(func(SuccessEvent)) EventListener
	OnFailureFunc(func(FailureEvent)) EventListener
	DismissSuccessFunc(func(SuccessEvent)) EventListener
	DismissFailureFunc(func(FailureEvent)) EventListener

	OnSuccess(ge.Action[SuccessEvent]) EventListener
	OnFailure(ge.Action[FailureEvent]) EventListener
	DismissSuccess(ge.Action[SuccessEvent]) EventListener
	DismissFailure(ge.Action[FailureEvent]) EventListener
}

func newEventListener() *eventListener {
	return &eventListener{
		onSuccess: make([]ge.Action[SuccessEvent], 0),
		onFailure: make([]ge.Action[FailureEvent], 0),
	}
}

type eventListener struct {
	sync.RWMutex
	onSuccess []ge.Action[SuccessEvent]
	onFailure []ge.Action[FailureEvent]
}

func (listener *eventListener) OnSuccessFunc(consumer func(SuccessEvent)) EventListener {
	return listener.OnSuccess(ge.ActionFunc[SuccessEvent](consumer))
}

func (listener *eventListener) OnFailureFunc(consumer func(FailureEvent)) EventListener {
	return listener.OnFailure(ge.ActionFunc[FailureEvent](consumer))
}

func (listener *eventListener) DismissSuccessFunc(consumer func(SuccessEvent)) EventListener {
	return listener.DismissSuccess(ge.ActionFunc[SuccessEvent](consumer))
}

func (listener *eventListener) DismissFailureFunc(consumer func(FailureEvent)) EventListener {
	return listener.DismissFailure(ge.ActionFunc[FailureEvent](consumer))
}

func (listener *eventListener) OnSuccess(action ge.Action[SuccessEvent]) EventListener {
	listener.Lock()
	defer listener.Unlock()
	listener.onSuccess = ge.AppendElementUnique(listener.onSuccess, action)
	return listener
}

func (listener *eventListener) OnFailure(action ge.Action[FailureEvent]) EventListener {
	listener.Lock()
	defer listener.Unlock()
	listener.onFailure = ge.AppendElementUnique(listener.onFailure, action)
	return listener
}

func (listener *eventListener) DismissSuccess(action ge.Action[SuccessEvent]) EventListener {
	listener.Lock()
	defer listener.Unlock()
	listener.onSuccess = ge.RemoveElementByValue(listener.onSuccess, action)
	return listener
}

func (listener *eventListener) DismissFailure(action ge.Action[FailureEvent]) EventListener {
	listener.Lock()
	defer listener.Unlock()
	listener.onFailure = ge.RemoveElementByValue(listener.onFailure, action)
	return listener
}

func (listener *eventListener) consumeEvent(event Event) {
	go func() {
		listener.RLock()
		defer listener.RUnlock()
		switch e := event.(type) {
		case *successEvent:
			ge.ForEach(listener.onSuccess, SuccessEvent(e))
		case *failureEvent:
			ge.ForEach(listener.onFailure, FailureEvent(e))
		}
	}()
}
