package circuitbreaker

import (
	"github.com/CharLemAznable/gogo/ext"
	. "github.com/CharLemAznable/gogo/fn"
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

	OnSuccess(Consumer[SuccessEvent]) EventListener
	OnError(Consumer[ErrorEvent]) EventListener
	OnNotPermitted(Consumer[NotPermittedEvent]) EventListener
	OnStateTransition(Consumer[StateTransitionEvent]) EventListener
	OnFailureRateExceeded(Consumer[FailureRateExceededEvent]) EventListener
	OnSlowCallRateExceeded(Consumer[SlowCallRateExceededEvent]) EventListener
	DismissSuccess(Consumer[SuccessEvent]) EventListener
	DismissError(Consumer[ErrorEvent]) EventListener
	DismissNotPermitted(Consumer[NotPermittedEvent]) EventListener
	DismissStateTransition(Consumer[StateTransitionEvent]) EventListener
	DismissFailureRateExceeded(Consumer[FailureRateExceededEvent]) EventListener
	DismissSlowCallRateExceeded(Consumer[SlowCallRateExceededEvent]) EventListener
}

func newEventListener() *eventListener {
	return &eventListener{
		onSuccess:              ext.NewConsumers[SuccessEvent](),
		onError:                ext.NewConsumers[ErrorEvent](),
		onNotPermitted:         ext.NewConsumers[NotPermittedEvent](),
		onStateTransition:      ext.NewConsumers[StateTransitionEvent](),
		onFailureRateExceeded:  ext.NewConsumers[FailureRateExceededEvent](),
		onSlowCallRateExceeded: ext.NewConsumers[SlowCallRateExceededEvent](),
	}
}

type eventListener struct {
	onSuccess              ext.Consumers[SuccessEvent]
	onError                ext.Consumers[ErrorEvent]
	onNotPermitted         ext.Consumers[NotPermittedEvent]
	onStateTransition      ext.Consumers[StateTransitionEvent]
	onFailureRateExceeded  ext.Consumers[FailureRateExceededEvent]
	onSlowCallRateExceeded ext.Consumers[SlowCallRateExceededEvent]
}

func (listener *eventListener) OnSuccessFunc(consumer func(SuccessEvent)) EventListener {
	return listener.OnSuccess(ConsumerOf(consumer))
}

func (listener *eventListener) OnErrorFunc(consumer func(ErrorEvent)) EventListener {
	return listener.OnError(ConsumerOf(consumer))
}

func (listener *eventListener) OnNotPermittedFunc(consumer func(NotPermittedEvent)) EventListener {
	return listener.OnNotPermitted(ConsumerOf(consumer))
}

func (listener *eventListener) OnStateTransitionFunc(consumer func(StateTransitionEvent)) EventListener {
	return listener.OnStateTransition(ConsumerOf(consumer))
}

func (listener *eventListener) OnFailureRateExceededFunc(consumer func(FailureRateExceededEvent)) EventListener {
	return listener.OnFailureRateExceeded(ConsumerOf(consumer))
}

func (listener *eventListener) OnSlowCallRateExceededFunc(consumer func(SlowCallRateExceededEvent)) EventListener {
	return listener.OnSlowCallRateExceeded(ConsumerOf(consumer))
}

func (listener *eventListener) DismissSuccessFunc(consumer func(SuccessEvent)) EventListener {
	return listener.DismissSuccess(ConsumerOf(consumer))
}

func (listener *eventListener) DismissErrorFunc(consumer func(ErrorEvent)) EventListener {
	return listener.DismissError(ConsumerOf(consumer))
}

func (listener *eventListener) DismissNotPermittedFunc(consumer func(NotPermittedEvent)) EventListener {
	return listener.DismissNotPermitted(ConsumerOf(consumer))
}

func (listener *eventListener) DismissStateTransitionFunc(consumer func(StateTransitionEvent)) EventListener {
	return listener.DismissStateTransition(ConsumerOf(consumer))
}

func (listener *eventListener) DismissFailureRateExceededFunc(consumer func(FailureRateExceededEvent)) EventListener {
	return listener.DismissFailureRateExceeded(ConsumerOf(consumer))
}

func (listener *eventListener) DismissSlowCallRateExceededFunc(consumer func(SlowCallRateExceededEvent)) EventListener {
	return listener.DismissSlowCallRateExceeded(ConsumerOf(consumer))
}

func (listener *eventListener) OnSuccess(consumer Consumer[SuccessEvent]) EventListener {
	listener.onSuccess.AppendConsumer(consumer)
	return listener
}

func (listener *eventListener) OnError(consumer Consumer[ErrorEvent]) EventListener {
	listener.onError.AppendConsumer(consumer)
	return listener
}

func (listener *eventListener) OnNotPermitted(consumer Consumer[NotPermittedEvent]) EventListener {
	listener.onNotPermitted.AppendConsumer(consumer)
	return listener
}

func (listener *eventListener) OnStateTransition(consumer Consumer[StateTransitionEvent]) EventListener {
	listener.onStateTransition.AppendConsumer(consumer)
	return listener
}

func (listener *eventListener) OnFailureRateExceeded(consumer Consumer[FailureRateExceededEvent]) EventListener {
	listener.onFailureRateExceeded.AppendConsumer(consumer)
	return listener
}

func (listener *eventListener) OnSlowCallRateExceeded(consumer Consumer[SlowCallRateExceededEvent]) EventListener {
	listener.onSlowCallRateExceeded.AppendConsumer(consumer)
	return listener
}

func (listener *eventListener) DismissSuccess(consumer Consumer[SuccessEvent]) EventListener {
	listener.onSuccess.RemoveConsumer(consumer)
	return listener
}

func (listener *eventListener) DismissError(consumer Consumer[ErrorEvent]) EventListener {
	listener.onError.RemoveConsumer(consumer)
	return listener
}

func (listener *eventListener) DismissNotPermitted(consumer Consumer[NotPermittedEvent]) EventListener {
	listener.onNotPermitted.RemoveConsumer(consumer)
	return listener
}

func (listener *eventListener) DismissStateTransition(consumer Consumer[StateTransitionEvent]) EventListener {
	listener.onStateTransition.RemoveConsumer(consumer)
	return listener
}

func (listener *eventListener) DismissFailureRateExceeded(consumer Consumer[FailureRateExceededEvent]) EventListener {
	listener.onFailureRateExceeded.RemoveConsumer(consumer)
	return listener
}

func (listener *eventListener) DismissSlowCallRateExceeded(consumer Consumer[SlowCallRateExceededEvent]) EventListener {
	listener.onSlowCallRateExceeded.RemoveConsumer(consumer)
	return listener
}

func (listener *eventListener) consumeEvent(event Event) {
	go func() {
		switch e := event.(type) {
		case *successEvent:
			listener.onSuccess.Accept(SuccessEvent(e))
		case *errorEvent:
			listener.onError.Accept(ErrorEvent(e))
		case *notPermittedEvent:
			listener.onNotPermitted.Accept(NotPermittedEvent(e))
		case *stateTransitionEvent:
			listener.onStateTransition.Accept(StateTransitionEvent(e))
		case *failureRateExceededEvent:
			listener.onFailureRateExceeded.Accept(FailureRateExceededEvent(e))
		case *slowCallRateExceededEvent:
			listener.onSlowCallRateExceeded.Accept(SlowCallRateExceededEvent(e))
		}
	}()
}
