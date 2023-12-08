package ratelimiter

import (
	"github.com/CharLemAznable/ge"
	"sync"
)

type EventListener interface {
	OnSuccess(func(SuccessEvent)) EventListener
	OnFailure(func(FailureEvent)) EventListener
	Dismiss(any) EventListener
}

func newEventListener() *eventListener {
	return &eventListener{
		onSuccess: make([]func(SuccessEvent), 0),
		onFailure: make([]func(FailureEvent), 0),
	}
}

type eventListener struct {
	sync.RWMutex
	onSuccess []func(SuccessEvent)
	onFailure []func(FailureEvent)
}

func (listener *eventListener) OnSuccess(consumer func(SuccessEvent)) EventListener {
	listener.Lock()
	defer listener.Unlock()
	listener.onSuccess = ge.AppendElementUnique(listener.onSuccess, consumer)
	return listener
}

func (listener *eventListener) OnFailure(consumer func(FailureEvent)) EventListener {
	listener.Lock()
	defer listener.Unlock()
	listener.onFailure = ge.AppendElementUnique(listener.onFailure, consumer)
	return listener
}

func (listener *eventListener) Dismiss(consumer any) EventListener {
	listener.Lock()
	defer listener.Unlock()
	switch c := consumer.(type) {
	case func(SuccessEvent):
		listener.onSuccess = ge.RemoveElementByValue(listener.onSuccess, c)
	case func(FailureEvent):
		listener.onFailure = ge.RemoveElementByValue(listener.onFailure, c)
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
		case *failureEvent:
			ge.ConsumeEach(listener.onFailure, FailureEvent(e))
		}
	}()
}
