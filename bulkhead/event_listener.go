package bulkhead

import (
	"github.com/CharLemAznable/resilience4go/utils"
	"sync"
)

type PermittedEventConsumer func(PermittedEvent)
type RejectedEventConsumer func(RejectedEvent)
type FinishedEventConsumer func(FinishedEvent)

type EventListener interface {
	OnPermitted(PermittedEventConsumer) EventListener
	OnRejected(RejectedEventConsumer) EventListener
	OnFinished(FinishedEventConsumer) EventListener
	Dismiss(any) EventListener
	consumeEvent(Event)
}

func newEventListener() EventListener {
	return &eventListener{
		onPermitted:       make([]PermittedEventConsumer, 0),
		onPermittedSlices: utils.NewSlicesWithPointer[PermittedEventConsumer](),
		onRejected:        make([]RejectedEventConsumer, 0),
		onRejectedSlices:  utils.NewSlicesWithPointer[RejectedEventConsumer](),
		onFinished:        make([]FinishedEventConsumer, 0),
		onFinishedSlices:  utils.NewSlicesWithPointer[FinishedEventConsumer](),
	}
}

type eventListener struct {
	sync.RWMutex
	onPermitted       []PermittedEventConsumer
	onPermittedSlices utils.Slices[PermittedEventConsumer]
	onRejected        []RejectedEventConsumer
	onRejectedSlices  utils.Slices[RejectedEventConsumer]
	onFinished        []FinishedEventConsumer
	onFinishedSlices  utils.Slices[FinishedEventConsumer]
}

func (listener *eventListener) OnPermitted(consumer PermittedEventConsumer) EventListener {
	listener.Lock()
	defer listener.Unlock()
	listener.onPermitted = listener.onPermittedSlices.AppendElementUnique(listener.onPermitted, consumer)
	return listener
}

func (listener *eventListener) OnRejected(consumer RejectedEventConsumer) EventListener {
	listener.Lock()
	defer listener.Unlock()
	listener.onRejected = listener.onRejectedSlices.AppendElementUnique(listener.onRejected, consumer)
	return listener
}

func (listener *eventListener) OnFinished(consumer FinishedEventConsumer) EventListener {
	listener.Lock()
	defer listener.Unlock()
	listener.onFinished = listener.onFinishedSlices.AppendElementUnique(listener.onFinished, consumer)
	return listener
}

func (listener *eventListener) Dismiss(consumer any) EventListener {
	listener.Lock()
	defer listener.Unlock()
	if c, ok := consumer.(func(PermittedEvent)); ok {
		listener.onPermitted = listener.onPermittedSlices.RemoveElementByValue(listener.onPermitted, c)
	}
	if c, ok := consumer.(func(RejectedEvent)); ok {
		listener.onRejected = listener.onRejectedSlices.RemoveElementByValue(listener.onRejected, c)
	}
	if c, ok := consumer.(func(FinishedEvent)); ok {
		listener.onFinished = listener.onFinishedSlices.RemoveElementByValue(listener.onFinished, c)
	}
	return listener
}

func (listener *eventListener) consumeEvent(event Event) {
	go func() {
		listener.RLock()
		defer listener.RUnlock()
		switch e := event.(type) {
		case *permittedEvent:
			for _, consumer := range listener.onPermitted {
				go consumer(e)
			}
		case *rejectedEvent:
			for _, consumer := range listener.onRejected {
				go consumer(e)
			}
		case *finishedEvent:
			for _, consumer := range listener.onFinished {
				go consumer(e)
			}
		}
	}()
}
