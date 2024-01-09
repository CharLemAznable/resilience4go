package retry

import (
	"github.com/CharLemAznable/gogo/ext"
	. "github.com/CharLemAznable/gogo/fn"
)

type EventListener interface {
	OnSuccessFunc(func(SuccessEvent)) EventListener
	OnRetryFunc(func(RetryEvent)) EventListener
	OnErrorFunc(func(ErrorEvent)) EventListener
	DismissSuccessFunc(func(SuccessEvent)) EventListener
	DismissRetryFunc(func(RetryEvent)) EventListener
	DismissErrorFunc(func(ErrorEvent)) EventListener

	OnSuccess(Consumer[SuccessEvent]) EventListener
	OnRetry(Consumer[RetryEvent]) EventListener
	OnError(Consumer[ErrorEvent]) EventListener
	DismissSuccess(Consumer[SuccessEvent]) EventListener
	DismissRetry(Consumer[RetryEvent]) EventListener
	DismissError(Consumer[ErrorEvent]) EventListener
}

func newEventListener() *eventListener {
	return &eventListener{
		onSuccess: ext.NewConsumers[SuccessEvent](),
		onRetry:   ext.NewConsumers[RetryEvent](),
		onError:   ext.NewConsumers[ErrorEvent](),
	}
}

type eventListener struct {
	onSuccess ext.Consumers[SuccessEvent]
	onRetry   ext.Consumers[RetryEvent]
	onError   ext.Consumers[ErrorEvent]
}

func (listener *eventListener) OnSuccessFunc(consumer func(SuccessEvent)) EventListener {
	return listener.OnSuccess(ConsumerOf(consumer))
}

func (listener *eventListener) OnRetryFunc(consumer func(RetryEvent)) EventListener {
	return listener.OnRetry(ConsumerOf(consumer))
}

func (listener *eventListener) OnErrorFunc(consumer func(ErrorEvent)) EventListener {
	return listener.OnError(ConsumerOf(consumer))
}

func (listener *eventListener) DismissSuccessFunc(consumer func(SuccessEvent)) EventListener {
	return listener.DismissSuccess(ConsumerOf(consumer))
}

func (listener *eventListener) DismissRetryFunc(consumer func(RetryEvent)) EventListener {
	return listener.DismissRetry(ConsumerOf(consumer))
}

func (listener *eventListener) DismissErrorFunc(consumer func(ErrorEvent)) EventListener {
	return listener.DismissError(ConsumerOf(consumer))
}

func (listener *eventListener) OnSuccess(consumer Consumer[SuccessEvent]) EventListener {
	listener.onSuccess.AppendConsumer(consumer)
	return listener
}

func (listener *eventListener) OnRetry(consumer Consumer[RetryEvent]) EventListener {
	listener.onRetry.AppendConsumer(consumer)
	return listener
}

func (listener *eventListener) OnError(consumer Consumer[ErrorEvent]) EventListener {
	listener.onError.AppendConsumer(consumer)
	return listener
}

func (listener *eventListener) DismissSuccess(consumer Consumer[SuccessEvent]) EventListener {
	listener.onSuccess.RemoveConsumer(consumer)
	return listener
}

func (listener *eventListener) DismissRetry(consumer Consumer[RetryEvent]) EventListener {
	listener.onRetry.RemoveConsumer(consumer)
	return listener
}

func (listener *eventListener) DismissError(consumer Consumer[ErrorEvent]) EventListener {
	listener.onError.RemoveConsumer(consumer)
	return listener
}

func (listener *eventListener) consumeEvent(event Event) {
	go func() {
		switch e := event.(type) {
		case *successEvent:
			listener.onSuccess.Accept(SuccessEvent(e))
		case *retryEvent:
			listener.onRetry.Accept(RetryEvent(e))
		case *errorEvent:
			listener.onError.Accept(ErrorEvent(e))
		}
	}()
}
