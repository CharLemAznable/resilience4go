package bulkhead

import (
	"github.com/CharLemAznable/gogo/ext"
	. "github.com/CharLemAznable/gogo/fn"
)

type EventListener interface {
	OnPermittedFunc(func(PermittedEvent)) EventListener
	OnRejectedFunc(func(RejectedEvent)) EventListener
	OnFinishedFunc(func(FinishedEvent)) EventListener
	DismissPermittedFunc(func(PermittedEvent)) EventListener
	DismissRejectedFunc(func(RejectedEvent)) EventListener
	DismissFinishedFunc(func(FinishedEvent)) EventListener

	OnPermitted(Consumer[PermittedEvent]) EventListener
	OnRejected(Consumer[RejectedEvent]) EventListener
	OnFinished(Consumer[FinishedEvent]) EventListener
	DismissPermitted(Consumer[PermittedEvent]) EventListener
	DismissRejected(Consumer[RejectedEvent]) EventListener
	DismissFinished(Consumer[FinishedEvent]) EventListener
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
	return listener.OnPermitted(ConsumerOf(consumer))
}

func (listener *eventListener) OnRejectedFunc(consumer func(RejectedEvent)) EventListener {
	return listener.OnRejected(ConsumerOf(consumer))
}

func (listener *eventListener) OnFinishedFunc(consumer func(FinishedEvent)) EventListener {
	return listener.OnFinished(ConsumerOf(consumer))
}

func (listener *eventListener) DismissPermittedFunc(consumer func(PermittedEvent)) EventListener {
	return listener.DismissPermitted(ConsumerOf(consumer))
}

func (listener *eventListener) DismissRejectedFunc(consumer func(RejectedEvent)) EventListener {
	return listener.DismissRejected(ConsumerOf(consumer))
}

func (listener *eventListener) DismissFinishedFunc(consumer func(FinishedEvent)) EventListener {
	return listener.DismissFinished(ConsumerOf(consumer))
}

func (listener *eventListener) OnPermitted(consumer Consumer[PermittedEvent]) EventListener {
	listener.onPermitted.AppendConsumer(consumer)
	return listener
}

func (listener *eventListener) OnRejected(consumer Consumer[RejectedEvent]) EventListener {
	listener.onRejected.AppendConsumer(consumer)
	return listener
}

func (listener *eventListener) OnFinished(consumer Consumer[FinishedEvent]) EventListener {
	listener.onFinished.AppendConsumer(consumer)
	return listener
}

func (listener *eventListener) DismissPermitted(consumer Consumer[PermittedEvent]) EventListener {
	listener.onPermitted.RemoveConsumer(consumer)
	return listener
}

func (listener *eventListener) DismissRejected(consumer Consumer[RejectedEvent]) EventListener {
	listener.onRejected.RemoveConsumer(consumer)
	return listener
}

func (listener *eventListener) DismissFinished(consumer Consumer[FinishedEvent]) EventListener {
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
