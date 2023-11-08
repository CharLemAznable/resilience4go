package bulkhead

import (
	"github.com/CharLemAznable/resilience4go/utils"
)

type EventConsumer func(Event)

type EventListener interface {
	OnPermitted(EventConsumer) EventListener
	OnRejected(EventConsumer) EventListener
	OnFinished(EventConsumer) EventListener
	Dismiss(EventConsumer) EventListener
	HasConsumer() bool
	consumeEvent(Event)
}

func newEventListener() EventListener {
	return &eventListener{
		onPermitted: make([]EventConsumer, 0),
		onRejected:  make([]EventConsumer, 0),
		onFinished:  make([]EventConsumer, 0),
		slices:      utils.NewSlicesWithPointer[EventConsumer](),
	}
}

type eventListener struct {
	onPermitted []EventConsumer
	onRejected  []EventConsumer
	onFinished  []EventConsumer
	slices      utils.Slices[EventConsumer]
}

func (listener *eventListener) OnPermitted(consumer EventConsumer) EventListener {
	listener.onPermitted = append(listener.onPermitted, consumer)
	return listener
}

func (listener *eventListener) OnRejected(consumer EventConsumer) EventListener {
	listener.onRejected = append(listener.onRejected, consumer)
	return listener
}

func (listener *eventListener) OnFinished(consumer EventConsumer) EventListener {
	listener.onFinished = append(listener.onFinished, consumer)
	return listener
}

func (listener *eventListener) Dismiss(consumer EventConsumer) EventListener {
	listener.onPermitted = listener.slices.RemoveElementByValue(listener.onPermitted, consumer)
	listener.onRejected = listener.slices.RemoveElementByValue(listener.onRejected, consumer)
	listener.onFinished = listener.slices.RemoveElementByValue(listener.onFinished, consumer)
	return listener
}

func (listener *eventListener) HasConsumer() bool {
	return len(listener.onPermitted) > 0 || len(listener.onRejected) > 0 || len(listener.onFinished) > 0
}

func (listener *eventListener) consumeEvent(event Event) {
	if !listener.HasConsumer() {
		return
	}
	var consumers []EventConsumer
	switch event.EventType() {
	case PERMITTED:
		consumers = listener.onPermitted
	case REJECTED:
		consumers = listener.onRejected
	case FINISHED:
		consumers = listener.onFinished
	}
	for _, consumer := range consumers {
		go consumer(event)
	}
}
