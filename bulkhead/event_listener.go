package bulkhead

import (
	"github.com/CharLemAznable/gogo/ext"
	"github.com/CharLemAznable/gogo/fn"
)

type EventListener interface {
	OnPermittedFunc(func(PermittedEvent)) EventListener
	OnRejectedFunc(func(RejectedEvent)) EventListener
	OnFinishedFunc(func(FinishedEvent)) EventListener
	DismissPermittedFunc(func(PermittedEvent)) EventListener
	DismissRejectedFunc(func(RejectedEvent)) EventListener
	DismissFinishedFunc(func(FinishedEvent)) EventListener

	OnPermitted(fn.Consumer[PermittedEvent]) EventListener
	OnRejected(fn.Consumer[RejectedEvent]) EventListener
	OnFinished(fn.Consumer[FinishedEvent]) EventListener
	DismissPermitted(fn.Consumer[PermittedEvent]) EventListener
	DismissRejected(fn.Consumer[RejectedEvent]) EventListener
	DismissFinished(fn.Consumer[FinishedEvent]) EventListener
}

func newEventListener() *eventListener {
	return &eventListener{
		onPermitted: ext.NewConsumers[PermittedEvent](),
		onRejected:  ext.NewConsumers[RejectedEvent](),
		onFinished:  ext.NewConsumers[FinishedEvent](),
	}
}

type eventListener struct {
	onPermitted ext.Consumers[PermittedEvent]
	onRejected  ext.Consumers[RejectedEvent]
	onFinished  ext.Consumers[FinishedEvent]
}

func (listener *eventListener) OnPermittedFunc(consumer func(PermittedEvent)) EventListener {
	return listener.OnPermitted(fn.ConsumerOf(consumer))
}

func (listener *eventListener) OnRejectedFunc(consumer func(RejectedEvent)) EventListener {
	return listener.OnRejected(fn.ConsumerOf(consumer))
}

func (listener *eventListener) OnFinishedFunc(consumer func(FinishedEvent)) EventListener {
	return listener.OnFinished(fn.ConsumerOf(consumer))
}

func (listener *eventListener) DismissPermittedFunc(consumer func(PermittedEvent)) EventListener {
	return listener.DismissPermitted(fn.ConsumerOf(consumer))
}

func (listener *eventListener) DismissRejectedFunc(consumer func(RejectedEvent)) EventListener {
	return listener.DismissRejected(fn.ConsumerOf(consumer))
}

func (listener *eventListener) DismissFinishedFunc(consumer func(FinishedEvent)) EventListener {
	return listener.DismissFinished(fn.ConsumerOf(consumer))
}

func (listener *eventListener) OnPermitted(consumer fn.Consumer[PermittedEvent]) EventListener {
	listener.onPermitted.AppendConsumer(consumer)
	return listener
}

func (listener *eventListener) OnRejected(consumer fn.Consumer[RejectedEvent]) EventListener {
	listener.onRejected.AppendConsumer(consumer)
	return listener
}

func (listener *eventListener) OnFinished(consumer fn.Consumer[FinishedEvent]) EventListener {
	listener.onFinished.AppendConsumer(consumer)
	return listener
}

func (listener *eventListener) DismissPermitted(consumer fn.Consumer[PermittedEvent]) EventListener {
	listener.onPermitted.RemoveConsumer(consumer)
	return listener
}

func (listener *eventListener) DismissRejected(consumer fn.Consumer[RejectedEvent]) EventListener {
	listener.onRejected.RemoveConsumer(consumer)
	return listener
}

func (listener *eventListener) DismissFinished(consumer fn.Consumer[FinishedEvent]) EventListener {
	listener.onFinished.RemoveConsumer(consumer)
	return listener
}

func (listener *eventListener) consumeEvent(event Event) {
	go func() {
		switch e := event.(type) {
		case *permittedEvent:
			listener.onPermitted.Accept(PermittedEvent(e))
		case *rejectedEvent:
			listener.onRejected.Accept(RejectedEvent(e))
		case *finishedEvent:
			listener.onFinished.Accept(FinishedEvent(e))
		}
	}()
}
