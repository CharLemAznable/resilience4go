package ratelimiter

import (
	"github.com/CharLemAznable/gogo/ext"
	"github.com/CharLemAznable/gogo/fn"
)

type EventListener interface {
	OnSuccessFunc(func(SuccessEvent)) EventListener
	OnFailureFunc(func(FailureEvent)) EventListener
	DismissSuccessFunc(func(SuccessEvent)) EventListener
	DismissFailureFunc(func(FailureEvent)) EventListener

	OnSuccess(fn.Consumer[SuccessEvent]) EventListener
	OnFailure(fn.Consumer[FailureEvent]) EventListener
	DismissSuccess(fn.Consumer[SuccessEvent]) EventListener
	DismissFailure(fn.Consumer[FailureEvent]) EventListener
}

func newEventListener() *eventListener {
	return &eventListener{
		onSuccess: ext.NewConsumers[SuccessEvent](),
		onFailure: ext.NewConsumers[FailureEvent](),
	}
}

type eventListener struct {
	onSuccess ext.Consumers[SuccessEvent]
	onFailure ext.Consumers[FailureEvent]
}

func (listener *eventListener) OnSuccessFunc(consumer func(SuccessEvent)) EventListener {
	return listener.OnSuccess(fn.ConsumerOf(consumer))
}

func (listener *eventListener) OnFailureFunc(consumer func(FailureEvent)) EventListener {
	return listener.OnFailure(fn.ConsumerOf(consumer))
}

func (listener *eventListener) DismissSuccessFunc(consumer func(SuccessEvent)) EventListener {
	return listener.DismissSuccess(fn.ConsumerOf(consumer))
}

func (listener *eventListener) DismissFailureFunc(consumer func(FailureEvent)) EventListener {
	return listener.DismissFailure(fn.ConsumerOf(consumer))
}

func (listener *eventListener) OnSuccess(consumer fn.Consumer[SuccessEvent]) EventListener {
	listener.onSuccess.AppendConsumer(consumer)
	return listener
}

func (listener *eventListener) OnFailure(consumer fn.Consumer[FailureEvent]) EventListener {
	listener.onFailure.AppendConsumer(consumer)
	return listener
}

func (listener *eventListener) DismissSuccess(consumer fn.Consumer[SuccessEvent]) EventListener {
	listener.onSuccess.RemoveConsumer(consumer)
	return listener
}

func (listener *eventListener) DismissFailure(consumer fn.Consumer[FailureEvent]) EventListener {
	listener.onFailure.RemoveConsumer(consumer)
	return listener
}

func (listener *eventListener) consumeEvent(event Event) {
	go func() {
		switch e := event.(type) {
		case *successEvent:
			listener.onSuccess.Accept(SuccessEvent(e))
		case *failureEvent:
			listener.onFailure.Accept(FailureEvent(e))
		}
	}()
}
