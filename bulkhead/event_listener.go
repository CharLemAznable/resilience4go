package bulkhead

import (
	"github.com/CharLemAznable/ge"
	"sync"
)

type EventListener interface {
	OnPermittedFunc(func(PermittedEvent)) EventListener
	OnRejectedFunc(func(RejectedEvent)) EventListener
	OnFinishedFunc(func(FinishedEvent)) EventListener
	DismissPermittedFunc(func(PermittedEvent)) EventListener
	DismissRejectedFunc(func(RejectedEvent)) EventListener
	DismissFinishedFunc(func(FinishedEvent)) EventListener

	OnPermitted(ge.Action[PermittedEvent]) EventListener
	OnRejected(ge.Action[RejectedEvent]) EventListener
	OnFinished(ge.Action[FinishedEvent]) EventListener
	DismissPermitted(ge.Action[PermittedEvent]) EventListener
	DismissRejected(ge.Action[RejectedEvent]) EventListener
	DismissFinished(ge.Action[FinishedEvent]) EventListener
}

func newEventListener() *eventListener {
	return &eventListener{
		onPermitted: make([]ge.Action[PermittedEvent], 0),
		onRejected:  make([]ge.Action[RejectedEvent], 0),
		onFinished:  make([]ge.Action[FinishedEvent], 0),
	}
}

type eventListener struct {
	sync.RWMutex
	onPermitted []ge.Action[PermittedEvent]
	onRejected  []ge.Action[RejectedEvent]
	onFinished  []ge.Action[FinishedEvent]
}

func (listener *eventListener) OnPermittedFunc(consumer func(PermittedEvent)) EventListener {
	return listener.OnPermitted(ge.ActionFunc[PermittedEvent](consumer))
}

func (listener *eventListener) OnRejectedFunc(consumer func(RejectedEvent)) EventListener {
	return listener.OnRejected(ge.ActionFunc[RejectedEvent](consumer))
}

func (listener *eventListener) OnFinishedFunc(consumer func(FinishedEvent)) EventListener {
	return listener.OnFinished(ge.ActionFunc[FinishedEvent](consumer))
}

func (listener *eventListener) DismissPermittedFunc(consumer func(PermittedEvent)) EventListener {
	return listener.DismissPermitted(ge.ActionFunc[PermittedEvent](consumer))
}

func (listener *eventListener) DismissRejectedFunc(consumer func(RejectedEvent)) EventListener {
	return listener.DismissRejected(ge.ActionFunc[RejectedEvent](consumer))
}

func (listener *eventListener) DismissFinishedFunc(consumer func(FinishedEvent)) EventListener {
	return listener.DismissFinished(ge.ActionFunc[FinishedEvent](consumer))
}

func (listener *eventListener) OnPermitted(action ge.Action[PermittedEvent]) EventListener {
	listener.Lock()
	defer listener.Unlock()
	listener.onPermitted = ge.AppendElementUnique(listener.onPermitted, action)
	return listener
}

func (listener *eventListener) OnRejected(action ge.Action[RejectedEvent]) EventListener {
	listener.Lock()
	defer listener.Unlock()
	listener.onRejected = ge.AppendElementUnique(listener.onRejected, action)
	return listener
}

func (listener *eventListener) OnFinished(action ge.Action[FinishedEvent]) EventListener {
	listener.Lock()
	defer listener.Unlock()
	listener.onFinished = ge.AppendElementUnique(listener.onFinished, action)
	return listener
}

func (listener *eventListener) DismissPermitted(action ge.Action[PermittedEvent]) EventListener {
	listener.Lock()
	defer listener.Unlock()
	listener.onPermitted = ge.RemoveElementByValue(listener.onPermitted, action)
	return listener
}

func (listener *eventListener) DismissRejected(action ge.Action[RejectedEvent]) EventListener {
	listener.Lock()
	defer listener.Unlock()
	listener.onRejected = ge.RemoveElementByValue(listener.onRejected, action)
	return listener
}

func (listener *eventListener) DismissFinished(action ge.Action[FinishedEvent]) EventListener {
	listener.Lock()
	defer listener.Unlock()
	listener.onFinished = ge.RemoveElementByValue(listener.onFinished, action)
	return listener
}

func (listener *eventListener) consumeEvent(event Event) {
	go func() {
		listener.RLock()
		defer listener.RUnlock()
		switch e := event.(type) {
		case *permittedEvent:
			ge.ForEach(listener.onPermitted, PermittedEvent(e))
		case *rejectedEvent:
			ge.ForEach(listener.onRejected, RejectedEvent(e))
		case *finishedEvent:
			ge.ForEach(listener.onFinished, FinishedEvent(e))
		}
	}()
}
