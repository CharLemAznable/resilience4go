package timelimiter

import (
	"github.com/CharLemAznable/resilience4go/utils"
	"sync"
)

type EventListener interface {
	OnSuccess(func(SuccessEvent)) EventListener
	OnTimeout(func(TimeoutEvent)) EventListener
	OnPanic(func(PanicEvent)) EventListener
	Dismiss(any) EventListener
	consumeEvent(Event)
}

func newEventListener() EventListener {
	return &eventListener{
		onSuccess: make([]func(SuccessEvent), 0),
		onTimeout: make([]func(TimeoutEvent), 0),
		onPanic:   make([]func(PanicEvent), 0),
	}
}

type eventListener struct {
	sync.RWMutex
	onSuccess []func(SuccessEvent)
	onTimeout []func(TimeoutEvent)
	onPanic   []func(PanicEvent)
}

func (listener *eventListener) OnSuccess(consumer func(SuccessEvent)) EventListener {
	listener.Lock()
	defer listener.Unlock()
	listener.onSuccess = utils.AppendElementUnique(listener.onSuccess, consumer)
	return listener
}

func (listener *eventListener) OnTimeout(consumer func(TimeoutEvent)) EventListener {
	listener.Lock()
	defer listener.Unlock()
	listener.onTimeout = utils.AppendElementUnique(listener.onTimeout, consumer)
	return listener
}

func (listener *eventListener) OnPanic(consumer func(PanicEvent)) EventListener {
	listener.Lock()
	defer listener.Unlock()
	listener.onPanic = utils.AppendElementUnique(listener.onPanic, consumer)
	return listener
}

func (listener *eventListener) Dismiss(consumer any) EventListener {
	listener.Lock()
	defer listener.Unlock()
	switch c := consumer.(type) {
	case func(SuccessEvent):
		listener.onSuccess = utils.RemoveElementByValue(listener.onSuccess, c)
	case func(TimeoutEvent):
		listener.onTimeout = utils.RemoveElementByValue(listener.onTimeout, c)
	case func(PanicEvent):
		listener.onPanic = utils.RemoveElementByValue(listener.onPanic, c)
	}
	return listener
}

func (listener *eventListener) consumeEvent(event Event) {
	go func() {
		listener.RLock()
		defer listener.RUnlock()
		switch e := event.(type) {
		case *successEvent:
			utils.ConsumeEvent(listener.onSuccess, SuccessEvent(e))
		case *timeoutEvent:
			utils.ConsumeEvent(listener.onTimeout, TimeoutEvent(e))
		case *panicEvent:
			utils.ConsumeEvent(listener.onPanic, PanicEvent(e))
		}
	}()
}
