package circuitbreaker

import (
	"github.com/CharLemAznable/ge"
	"sync"
)

type EventListener interface {
	OnSuccessFunc(func(SuccessEvent)) EventListener
	OnErrorFunc(func(ErrorEvent)) EventListener
	OnNotPermittedFunc(func(NotPermittedEvent)) EventListener
	OnStateTransitionFunc(func(StateTransitionEvent)) EventListener
	OnFailureRateExceededFunc(func(FailureRateExceededEvent)) EventListener
	OnSlowCallRateExceededFunc(func(SlowCallRateExceededEvent)) EventListener
	DismissSuccessFunc(func(SuccessEvent)) EventListener
	DismissErrorFunc(func(ErrorEvent)) EventListener
	DismissNotPermittedFunc(func(NotPermittedEvent)) EventListener
	DismissStateTransitionFunc(func(StateTransitionEvent)) EventListener
	DismissFailureRateExceededFunc(func(FailureRateExceededEvent)) EventListener
	DismissSlowCallRateExceededFunc(func(SlowCallRateExceededEvent)) EventListener

	OnSuccess(ge.Action[SuccessEvent]) EventListener
	OnError(ge.Action[ErrorEvent]) EventListener
	OnNotPermitted(ge.Action[NotPermittedEvent]) EventListener
	OnStateTransition(ge.Action[StateTransitionEvent]) EventListener
	OnFailureRateExceeded(ge.Action[FailureRateExceededEvent]) EventListener
	OnSlowCallRateExceeded(ge.Action[SlowCallRateExceededEvent]) EventListener
	DismissSuccess(ge.Action[SuccessEvent]) EventListener
	DismissError(ge.Action[ErrorEvent]) EventListener
	DismissNotPermitted(ge.Action[NotPermittedEvent]) EventListener
	DismissStateTransition(ge.Action[StateTransitionEvent]) EventListener
	DismissFailureRateExceeded(ge.Action[FailureRateExceededEvent]) EventListener
	DismissSlowCallRateExceeded(ge.Action[SlowCallRateExceededEvent]) EventListener
}

func newEventListener() *eventListener {
	return &eventListener{
		onSuccess:              make([]ge.Action[SuccessEvent], 0),
		onError:                make([]ge.Action[ErrorEvent], 0),
		onNotPermitted:         make([]ge.Action[NotPermittedEvent], 0),
		onStateTransition:      make([]ge.Action[StateTransitionEvent], 0),
		onFailureRateExceeded:  make([]ge.Action[FailureRateExceededEvent], 0),
		onSlowCallRateExceeded: make([]ge.Action[SlowCallRateExceededEvent], 0),
	}
}

type eventListener struct {
	sync.RWMutex
	onSuccess              []ge.Action[SuccessEvent]
	onError                []ge.Action[ErrorEvent]
	onNotPermitted         []ge.Action[NotPermittedEvent]
	onStateTransition      []ge.Action[StateTransitionEvent]
	onFailureRateExceeded  []ge.Action[FailureRateExceededEvent]
	onSlowCallRateExceeded []ge.Action[SlowCallRateExceededEvent]
}

func (listener *eventListener) OnSuccessFunc(consumer func(SuccessEvent)) EventListener {
	return listener.OnSuccess(ge.ActionFunc[SuccessEvent](consumer))
}

func (listener *eventListener) OnErrorFunc(consumer func(ErrorEvent)) EventListener {
	return listener.OnError(ge.ActionFunc[ErrorEvent](consumer))
}

func (listener *eventListener) OnNotPermittedFunc(consumer func(NotPermittedEvent)) EventListener {
	return listener.OnNotPermitted(ge.ActionFunc[NotPermittedEvent](consumer))
}

func (listener *eventListener) OnStateTransitionFunc(consumer func(StateTransitionEvent)) EventListener {
	return listener.OnStateTransition(ge.ActionFunc[StateTransitionEvent](consumer))
}

func (listener *eventListener) OnFailureRateExceededFunc(consumer func(FailureRateExceededEvent)) EventListener {
	return listener.OnFailureRateExceeded(ge.ActionFunc[FailureRateExceededEvent](consumer))
}

func (listener *eventListener) OnSlowCallRateExceededFunc(consumer func(SlowCallRateExceededEvent)) EventListener {
	return listener.OnSlowCallRateExceeded(ge.ActionFunc[SlowCallRateExceededEvent](consumer))
}

func (listener *eventListener) DismissSuccessFunc(consumer func(SuccessEvent)) EventListener {
	return listener.DismissSuccess(ge.ActionFunc[SuccessEvent](consumer))
}

func (listener *eventListener) DismissErrorFunc(consumer func(ErrorEvent)) EventListener {
	return listener.DismissError(ge.ActionFunc[ErrorEvent](consumer))
}

func (listener *eventListener) DismissNotPermittedFunc(consumer func(NotPermittedEvent)) EventListener {
	return listener.DismissNotPermitted(ge.ActionFunc[NotPermittedEvent](consumer))
}

func (listener *eventListener) DismissStateTransitionFunc(consumer func(StateTransitionEvent)) EventListener {
	return listener.DismissStateTransition(ge.ActionFunc[StateTransitionEvent](consumer))
}

func (listener *eventListener) DismissFailureRateExceededFunc(consumer func(FailureRateExceededEvent)) EventListener {
	return listener.DismissFailureRateExceeded(ge.ActionFunc[FailureRateExceededEvent](consumer))
}

func (listener *eventListener) DismissSlowCallRateExceededFunc(consumer func(SlowCallRateExceededEvent)) EventListener {
	return listener.DismissSlowCallRateExceeded(ge.ActionFunc[SlowCallRateExceededEvent](consumer))
}

func (listener *eventListener) OnSuccess(action ge.Action[SuccessEvent]) EventListener {
	listener.Lock()
	defer listener.Unlock()
	listener.onSuccess = ge.AppendElementUnique(listener.onSuccess, action)
	return listener
}

func (listener *eventListener) OnError(action ge.Action[ErrorEvent]) EventListener {
	listener.Lock()
	defer listener.Unlock()
	listener.onError = ge.AppendElementUnique(listener.onError, action)
	return listener
}

func (listener *eventListener) OnNotPermitted(action ge.Action[NotPermittedEvent]) EventListener {
	listener.Lock()
	defer listener.Unlock()
	listener.onNotPermitted = ge.AppendElementUnique(listener.onNotPermitted, action)
	return listener
}

func (listener *eventListener) OnStateTransition(action ge.Action[StateTransitionEvent]) EventListener {
	listener.Lock()
	defer listener.Unlock()
	listener.onStateTransition = ge.AppendElementUnique(listener.onStateTransition, action)
	return listener
}

func (listener *eventListener) OnFailureRateExceeded(action ge.Action[FailureRateExceededEvent]) EventListener {
	listener.Lock()
	defer listener.Unlock()
	listener.onFailureRateExceeded = ge.AppendElementUnique(listener.onFailureRateExceeded, action)
	return listener
}

func (listener *eventListener) OnSlowCallRateExceeded(action ge.Action[SlowCallRateExceededEvent]) EventListener {
	listener.Lock()
	defer listener.Unlock()
	listener.onSlowCallRateExceeded = ge.AppendElementUnique(listener.onSlowCallRateExceeded, action)
	return listener
}

func (listener *eventListener) DismissSuccess(action ge.Action[SuccessEvent]) EventListener {
	listener.Lock()
	defer listener.Unlock()
	listener.onSuccess = ge.RemoveElementByValue(listener.onSuccess, action)
	return listener
}

func (listener *eventListener) DismissError(action ge.Action[ErrorEvent]) EventListener {
	listener.Lock()
	defer listener.Unlock()
	listener.onError = ge.RemoveElementByValue(listener.onError, action)
	return listener
}

func (listener *eventListener) DismissNotPermitted(action ge.Action[NotPermittedEvent]) EventListener {
	listener.Lock()
	defer listener.Unlock()
	listener.onNotPermitted = ge.RemoveElementByValue(listener.onNotPermitted, action)
	return listener
}

func (listener *eventListener) DismissStateTransition(action ge.Action[StateTransitionEvent]) EventListener {
	listener.Lock()
	defer listener.Unlock()
	listener.onStateTransition = ge.RemoveElementByValue(listener.onStateTransition, action)
	return listener
}

func (listener *eventListener) DismissFailureRateExceeded(action ge.Action[FailureRateExceededEvent]) EventListener {
	listener.Lock()
	defer listener.Unlock()
	listener.onFailureRateExceeded = ge.RemoveElementByValue(listener.onFailureRateExceeded, action)
	return listener
}

func (listener *eventListener) DismissSlowCallRateExceeded(action ge.Action[SlowCallRateExceededEvent]) EventListener {
	listener.Lock()
	defer listener.Unlock()
	listener.onSlowCallRateExceeded = ge.RemoveElementByValue(listener.onSlowCallRateExceeded, action)
	return listener
}

func (listener *eventListener) consumeEvent(event Event) {
	go func() {
		listener.RLock()
		defer listener.RUnlock()
		switch e := event.(type) {
		case *successEvent:
			ge.ForEach(listener.onSuccess, SuccessEvent(e))
		case *errorEvent:
			ge.ForEach(listener.onError, ErrorEvent(e))
		case *notPermittedEvent:
			ge.ForEach(listener.onNotPermitted, NotPermittedEvent(e))
		case *stateTransitionEvent:
			ge.ForEach(listener.onStateTransition, StateTransitionEvent(e))
		case *failureRateExceededEvent:
			ge.ForEach(listener.onFailureRateExceeded, FailureRateExceededEvent(e))
		case *slowCallRateExceededEvent:
			ge.ForEach(listener.onSlowCallRateExceeded, SlowCallRateExceededEvent(e))
		}
	}()
}
