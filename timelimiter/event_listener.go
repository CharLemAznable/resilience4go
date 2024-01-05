package timelimiter

import (
	"github.com/CharLemAznable/gogo/ext"
	"github.com/CharLemAznable/gogo/fn"
)

type EventListener interface {
	OnSuccessFunc(func(SuccessEvent)) EventListener
	OnTimeoutFunc(func(TimeoutEvent)) EventListener
	OnPanicFunc(func(PanicEvent)) EventListener
	DismissSuccessFunc(func(SuccessEvent)) EventListener
	DismissTimeoutFunc(func(TimeoutEvent)) EventListener
	DismissPanicFunc(func(PanicEvent)) EventListener

	OnSuccess(fn.Consumer[SuccessEvent]) EventListener
	OnTimeout(fn.Consumer[TimeoutEvent]) EventListener
	OnPanic(fn.Consumer[PanicEvent]) EventListener
	DismissSuccess(fn.Consumer[SuccessEvent]) EventListener
	DismissTimeout(fn.Consumer[TimeoutEvent]) EventListener
	DismissPanic(fn.Consumer[PanicEvent]) EventListener
}

func newEventListener() *eventListener {
	return &eventListener{
		onSuccess: ext.NewConsumers[SuccessEvent](),
		onTimeout: ext.NewConsumers[TimeoutEvent](),
		onPanic:   ext.NewConsumers[PanicEvent](),
	}
}

type eventListener struct {
	onSuccess ext.Consumers[SuccessEvent]
	onTimeout ext.Consumers[TimeoutEvent]
	onPanic   ext.Consumers[PanicEvent]
}

func (listener *eventListener) OnSuccessFunc(consumer func(SuccessEvent)) EventListener {
	return listener.OnSuccess(fn.ConsumerOf(consumer))
}

func (listener *eventListener) OnTimeoutFunc(consumer func(TimeoutEvent)) EventListener {
	return listener.OnTimeout(fn.ConsumerOf(consumer))
}

func (listener *eventListener) OnPanicFunc(consumer func(PanicEvent)) EventListener {
	return listener.OnPanic(fn.ConsumerOf(consumer))
}

func (listener *eventListener) DismissSuccessFunc(consumer func(SuccessEvent)) EventListener {
	return listener.DismissSuccess(fn.ConsumerOf(consumer))
}

func (listener *eventListener) DismissTimeoutFunc(consumer func(TimeoutEvent)) EventListener {
	return listener.DismissTimeout(fn.ConsumerOf(consumer))
}

func (listener *eventListener) DismissPanicFunc(consumer func(PanicEvent)) EventListener {
	return listener.DismissPanic(fn.ConsumerOf(consumer))
}

func (listener *eventListener) OnSuccess(consumer fn.Consumer[SuccessEvent]) EventListener {
	listener.onSuccess.AppendConsumer(consumer)
	return listener
}

func (listener *eventListener) OnTimeout(consumer fn.Consumer[TimeoutEvent]) EventListener {
	listener.onTimeout.AppendConsumer(consumer)
	return listener
}

func (listener *eventListener) OnPanic(consumer fn.Consumer[PanicEvent]) EventListener {
	listener.onPanic.AppendConsumer(consumer)
	return listener
}

func (listener *eventListener) DismissSuccess(consumer fn.Consumer[SuccessEvent]) EventListener {
	listener.onSuccess.RemoveConsumer(consumer)
	return listener
}

func (listener *eventListener) DismissTimeout(consumer fn.Consumer[TimeoutEvent]) EventListener {
	listener.onTimeout.RemoveConsumer(consumer)
	return listener
}

func (listener *eventListener) DismissPanic(consumer fn.Consumer[PanicEvent]) EventListener {
	listener.onPanic.RemoveConsumer(consumer)
	return listener
}

func (listener *eventListener) consumeEvent(event Event) {
	go func() {
		switch e := event.(type) {
		case *successEvent:
			listener.onSuccess.Accept(SuccessEvent(e))
		case *timeoutEvent:
			listener.onTimeout.Accept(TimeoutEvent(e))
		case *panicEvent:
			listener.onPanic.Accept(PanicEvent(e))
		}
	}()
}
