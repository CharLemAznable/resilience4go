package circuitbreaker

import (
	"github.com/CharLemAznable/ge"
	"sync"
)

type EventListener interface {
	OnSuccess(func(SuccessEvent)) EventListener
	OnError(func(ErrorEvent)) EventListener
	OnNotPermitted(func(NotPermittedEvent)) EventListener
	OnStateTransition(func(StateTransitionEvent)) EventListener
	OnFailureRateExceeded(func(FailureRateExceededEvent)) EventListener
	OnSlowCallRateExceeded(func(SlowCallRateExceededEvent)) EventListener
	Dismiss(any) EventListener
}

func newEventListener() *eventListener {
	return &eventListener{
		onSuccess:              make([]func(SuccessEvent), 0),
		onError:                make([]func(ErrorEvent), 0),
		onNotPermitted:         make([]func(NotPermittedEvent), 0),
		onStateTransition:      make([]func(StateTransitionEvent), 0),
		onFailureRateExceeded:  make([]func(FailureRateExceededEvent), 0),
		onSlowCallRateExceeded: make([]func(SlowCallRateExceededEvent), 0),
	}
}

type eventListener struct {
	sync.RWMutex
	onSuccess              []func(SuccessEvent)
	onError                []func(ErrorEvent)
	onNotPermitted         []func(NotPermittedEvent)
	onStateTransition      []func(StateTransitionEvent)
	onFailureRateExceeded  []func(FailureRateExceededEvent)
	onSlowCallRateExceeded []func(SlowCallRateExceededEvent)
}

func (listener *eventListener) OnSuccess(consumer func(SuccessEvent)) EventListener {
	listener.Lock()
	defer listener.Unlock()
	listener.onSuccess = ge.AppendElementUnique(listener.onSuccess, consumer)
	return listener
}

func (listener *eventListener) OnError(consumer func(ErrorEvent)) EventListener {
	listener.Lock()
	defer listener.Unlock()
	listener.onError = ge.AppendElementUnique(listener.onError, consumer)
	return listener
}

func (listener *eventListener) OnNotPermitted(consumer func(NotPermittedEvent)) EventListener {
	listener.Lock()
	defer listener.Unlock()
	listener.onNotPermitted = ge.AppendElementUnique(listener.onNotPermitted, consumer)
	return listener
}

func (listener *eventListener) OnStateTransition(consumer func(StateTransitionEvent)) EventListener {
	listener.Lock()
	defer listener.Unlock()
	listener.onStateTransition = ge.AppendElementUnique(listener.onStateTransition, consumer)
	return listener
}

func (listener *eventListener) OnFailureRateExceeded(consumer func(FailureRateExceededEvent)) EventListener {
	listener.Lock()
	defer listener.Unlock()
	listener.onFailureRateExceeded = ge.AppendElementUnique(listener.onFailureRateExceeded, consumer)
	return listener
}

func (listener *eventListener) OnSlowCallRateExceeded(consumer func(SlowCallRateExceededEvent)) EventListener {
	listener.Lock()
	defer listener.Unlock()
	listener.onSlowCallRateExceeded = ge.AppendElementUnique(listener.onSlowCallRateExceeded, consumer)
	return listener
}

func (listener *eventListener) Dismiss(consumer any) EventListener {
	listener.Lock()
	defer listener.Unlock()
	switch c := consumer.(type) {
	case func(SuccessEvent):
		listener.onSuccess = ge.RemoveElementByValue(listener.onSuccess, c)
	case func(ErrorEvent):
		listener.onError = ge.RemoveElementByValue(listener.onError, c)
	case func(NotPermittedEvent):
		listener.onNotPermitted = ge.RemoveElementByValue(listener.onNotPermitted, c)
	case func(StateTransitionEvent):
		listener.onStateTransition = ge.RemoveElementByValue(listener.onStateTransition, c)
	case func(FailureRateExceededEvent):
		listener.onFailureRateExceeded = ge.RemoveElementByValue(listener.onFailureRateExceeded, c)
	case func(SlowCallRateExceededEvent):
		listener.onSlowCallRateExceeded = ge.RemoveElementByValue(listener.onSlowCallRateExceeded, c)
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
		case *errorEvent:
			ge.ConsumeEach(listener.onError, ErrorEvent(e))
		case *notPermittedEvent:
			ge.ConsumeEach(listener.onNotPermitted, NotPermittedEvent(e))
		case *stateTransitionEvent:
			ge.ConsumeEach(listener.onStateTransition, StateTransitionEvent(e))
		case *failureRateExceededEvent:
			ge.ConsumeEach(listener.onFailureRateExceeded, FailureRateExceededEvent(e))
		case *slowCallRateExceededEvent:
			ge.ConsumeEach(listener.onSlowCallRateExceeded, SlowCallRateExceededEvent(e))
		}
	}()
}
