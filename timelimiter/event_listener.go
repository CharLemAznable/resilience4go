package timelimiter

import (
	"github.com/CharLemAznable/gogo/ext"
	. "github.com/CharLemAznable/gogo/fn"
)

type EventListener interface {
	OnSuccessFunc(func(SuccessEvent)) EventListener
	OnErrorFunc(func(ErrorEvent)) EventListener
	OnTimeoutFunc(func(TimeoutEvent)) EventListener
	OnPanicFunc(func(PanicEvent)) EventListener
	DismissSuccessFunc(func(SuccessEvent)) EventListener
	DismissErrorFunc(func(ErrorEvent)) EventListener
	DismissTimeoutFunc(func(TimeoutEvent)) EventListener
	DismissPanicFunc(func(PanicEvent)) EventListener

	OnSuccess(Consumer[SuccessEvent]) EventListener
	OnError(Consumer[ErrorEvent]) EventListener
	OnTimeout(Consumer[TimeoutEvent]) EventListener
	OnPanic(Consumer[PanicEvent]) EventListener
	DismissSuccess(Consumer[SuccessEvent]) EventListener
	DismissError(Consumer[ErrorEvent]) EventListener
	DismissTimeout(Consumer[TimeoutEvent]) EventListener
	DismissPanic(Consumer[PanicEvent]) EventListener
}

func newEventListener() *eventListener {
	return &eventListener{
		onSuccess: ext.NewConsumers[SuccessEvent](),
		onError:   ext.NewConsumers[ErrorEvent](),
		onTimeout: ext.NewConsumers[TimeoutEvent](),
		onPanic:   ext.NewConsumers[PanicEvent](),
	}
}

type eventListener struct {
	onSuccess ext.Consumers[SuccessEvent]
	onError   ext.Consumers[ErrorEvent]
	onTimeout ext.Consumers[TimeoutEvent]
	onPanic   ext.Consumers[PanicEvent]
}

func (listener *eventListener) OnSuccessFunc(consumer func(SuccessEvent)) EventListener {
	return listener.OnSuccess(ConsumerOf(consumer))
}

func (listener *eventListener) OnErrorFunc(consumer func(ErrorEvent)) EventListener {
	return listener.OnError(ConsumerOf(consumer))
}

func (listener *eventListener) OnTimeoutFunc(consumer func(TimeoutEvent)) EventListener {
	return listener.OnTimeout(ConsumerOf(consumer))
}

func (listener *eventListener) OnPanicFunc(consumer func(PanicEvent)) EventListener {
	return listener.OnPanic(ConsumerOf(consumer))
}

func (listener *eventListener) DismissSuccessFunc(consumer func(SuccessEvent)) EventListener {
	return listener.DismissSuccess(ConsumerOf(consumer))
}

func (listener *eventListener) DismissErrorFunc(consumer func(ErrorEvent)) EventListener {
	return listener.DismissError(ConsumerOf(consumer))
}

func (listener *eventListener) DismissTimeoutFunc(consumer func(TimeoutEvent)) EventListener {
	return listener.DismissTimeout(ConsumerOf(consumer))
}

func (listener *eventListener) DismissPanicFunc(consumer func(PanicEvent)) EventListener {
	return listener.DismissPanic(ConsumerOf(consumer))
}

func (listener *eventListener) OnSuccess(consumer Consumer[SuccessEvent]) EventListener {
	listener.onSuccess.AppendConsumer(consumer)
	return listener
}

func (listener *eventListener) OnError(consumer Consumer[ErrorEvent]) EventListener {
	listener.onError.AppendConsumer(consumer)
	return listener
}

func (listener *eventListener) OnTimeout(consumer Consumer[TimeoutEvent]) EventListener {
	listener.onTimeout.AppendConsumer(consumer)
	return listener
}

func (listener *eventListener) OnPanic(consumer Consumer[PanicEvent]) EventListener {
	listener.onPanic.AppendConsumer(consumer)
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

func (listener *eventListener) DismissTimeout(consumer Consumer[TimeoutEvent]) EventListener {
	listener.onTimeout.RemoveConsumer(consumer)
	return listener
}

func (listener *eventListener) DismissPanic(consumer Consumer[PanicEvent]) EventListener {
	listener.onPanic.RemoveConsumer(consumer)
	return listener
}

func (listener *eventListener) consumeEvent(event Event) {
	go func() {
		switch e := event.(type) {
		case *successEvent:
			listener.onSuccess.Accept(SuccessEvent(e))
		case *errorEvent:
			listener.onError.Accept(ErrorEvent(e))
		case *timeoutEvent:
			listener.onTimeout.Accept(TimeoutEvent(e))
		case *panicEvent:
			listener.onPanic.Accept(PanicEvent(e))
		}
	}()
}
