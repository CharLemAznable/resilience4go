package retry

import "github.com/CharLemAznable/resilience4go/utils"

type EventConsumer func(Event)

type EventListener interface {
	OnSuccess(EventConsumer) EventListener
	OnError(EventConsumer) EventListener
	OnRetry(EventConsumer) EventListener
	Dismiss(EventConsumer) EventListener
	HasConsumer() bool
	consumeEvent(Event)
}

func newEventListener() EventListener {
	return &eventListener{
		onSuccess: make([]EventConsumer, 0),
		onError:   make([]EventConsumer, 0),
		onRetry:   make([]EventConsumer, 0),
		slices:    utils.NewSlicesWithPointer[EventConsumer](),
	}
}

type eventListener struct {
	onSuccess []EventConsumer
	onError   []EventConsumer
	onRetry   []EventConsumer
	slices    utils.Slices[EventConsumer]
}

func (listener *eventListener) OnSuccess(consumer EventConsumer) EventListener {
	listener.onSuccess = listener.slices.AppendElementUnique(listener.onSuccess, consumer)
	return listener
}

func (listener *eventListener) OnError(consumer EventConsumer) EventListener {
	listener.onError = listener.slices.AppendElementUnique(listener.onError, consumer)
	return listener
}

func (listener *eventListener) OnRetry(consumer EventConsumer) EventListener {
	listener.onRetry = listener.slices.AppendElementUnique(listener.onRetry, consumer)
	return listener
}

func (listener *eventListener) Dismiss(consumer EventConsumer) EventListener {
	listener.onSuccess = listener.slices.RemoveElementByValue(listener.onSuccess, consumer)
	listener.onError = listener.slices.RemoveElementByValue(listener.onError, consumer)
	listener.onRetry = listener.slices.RemoveElementByValue(listener.onRetry, consumer)
	return listener
}

func (listener *eventListener) HasConsumer() bool {
	return len(listener.onSuccess) > 0 || len(listener.onError) > 0 || len(listener.onRetry) > 0
}

func (listener *eventListener) consumeEvent(event Event) {
	if !listener.HasConsumer() {
		return
	}
	var consumers []EventConsumer
	switch event.EventType() {
	case SUCCESS:
		consumers = listener.onSuccess
	case ERROR:
		consumers = listener.onError
	case RETRY:
		consumers = listener.onRetry
	}
	for _, consumer := range consumers {
		go consumer(event)
	}
}
