package bulkhead

import (
	"github.com/CharLemAznable/ge"
	"sync"
)

type EventListener interface {
	OnPermitted(func(PermittedEvent)) EventListener
	OnRejected(func(RejectedEvent)) EventListener
	OnFinished(func(FinishedEvent)) EventListener
	Dismiss(any) EventListener
}

func newEventListener() *eventListener {
	return &eventListener{
		onPermitted: make([]func(PermittedEvent), 0),
		onRejected:  make([]func(RejectedEvent), 0),
		onFinished:  make([]func(FinishedEvent), 0),
	}
}

type eventListener struct {
	sync.RWMutex
	onPermitted []func(PermittedEvent)
	onRejected  []func(RejectedEvent)
	onFinished  []func(FinishedEvent)
}

func (listener *eventListener) OnPermitted(consumer func(PermittedEvent)) EventListener {
	listener.Lock()
	defer listener.Unlock()
	listener.onPermitted = ge.AppendElementUnique(listener.onPermitted, consumer)
	return listener
}

func (listener *eventListener) OnRejected(consumer func(RejectedEvent)) EventListener {
	listener.Lock()
	defer listener.Unlock()
	listener.onRejected = ge.AppendElementUnique(listener.onRejected, consumer)
	return listener
}

func (listener *eventListener) OnFinished(consumer func(FinishedEvent)) EventListener {
	listener.Lock()
	defer listener.Unlock()
	listener.onFinished = ge.AppendElementUnique(listener.onFinished, consumer)
	return listener
}

func (listener *eventListener) Dismiss(consumer any) EventListener {
	listener.Lock()
	defer listener.Unlock()
	switch c := consumer.(type) {
	case func(PermittedEvent):
		listener.onPermitted = ge.RemoveElementByValue(listener.onPermitted, c)
	case func(RejectedEvent):
		listener.onRejected = ge.RemoveElementByValue(listener.onRejected, c)
	case func(FinishedEvent):
		listener.onFinished = ge.RemoveElementByValue(listener.onFinished, c)
	}
	return listener
}

func (listener *eventListener) consumeEvent(event Event) {
	go func() {
		listener.RLock()
		defer listener.RUnlock()
		switch e := event.(type) {
		case *permittedEvent:
			ge.ConsumeEach(listener.onPermitted, PermittedEvent(e))
		case *rejectedEvent:
			ge.ConsumeEach(listener.onRejected, RejectedEvent(e))
		case *finishedEvent:
			ge.ConsumeEach(listener.onFinished, FinishedEvent(e))
		}
	}()
}
