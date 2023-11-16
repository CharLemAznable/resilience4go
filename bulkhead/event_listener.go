package bulkhead

import (
	"github.com/CharLemAznable/resilience4go/utils"
	"sync"
)

type EventListener interface {
	OnPermitted(func(PermittedEvent)) EventListener
	OnRejected(func(RejectedEvent)) EventListener
	OnFinished(func(FinishedEvent)) EventListener
	Dismiss(any) EventListener
	consumeEvent(Event)
}

func newEventListener() EventListener {
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
	listener.onPermitted = utils.AppendElementUnique(listener.onPermitted, consumer)
	return listener
}

func (listener *eventListener) OnRejected(consumer func(RejectedEvent)) EventListener {
	listener.Lock()
	defer listener.Unlock()
	listener.onRejected = utils.AppendElementUnique(listener.onRejected, consumer)
	return listener
}

func (listener *eventListener) OnFinished(consumer func(FinishedEvent)) EventListener {
	listener.Lock()
	defer listener.Unlock()
	listener.onFinished = utils.AppendElementUnique(listener.onFinished, consumer)
	return listener
}

func (listener *eventListener) Dismiss(consumer any) EventListener {
	listener.Lock()
	defer listener.Unlock()
	if c, ok := consumer.(func(PermittedEvent)); ok {
		listener.onPermitted = utils.RemoveElementByValue(listener.onPermitted, c)
	}
	if c, ok := consumer.(func(RejectedEvent)); ok {
		listener.onRejected = utils.RemoveElementByValue(listener.onRejected, c)
	}
	if c, ok := consumer.(func(FinishedEvent)); ok {
		listener.onFinished = utils.RemoveElementByValue(listener.onFinished, c)
	}
	return listener
}

func (listener *eventListener) consumeEvent(event Event) {
	go func() {
		listener.RLock()
		defer listener.RUnlock()
		switch e := event.(type) {
		case *permittedEvent:
			utils.ConsumeEvent(listener.onPermitted, PermittedEvent(e))
		case *rejectedEvent:
			utils.ConsumeEvent(listener.onRejected, RejectedEvent(e))
		case *finishedEvent:
			utils.ConsumeEvent(listener.onFinished, FinishedEvent(e))
		}
	}()
}
