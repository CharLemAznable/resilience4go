package timelimiter

import (
	"github.com/CharLemAznable/ge"
	"sync"
)

type EventListener interface {
	OnSuccessFunc(func(SuccessEvent)) EventListener
	OnTimeoutFunc(func(TimeoutEvent)) EventListener
	OnPanicFunc(func(PanicEvent)) EventListener
	DismissSuccessFunc(func(SuccessEvent)) EventListener
	DismissTimeoutFunc(func(TimeoutEvent)) EventListener
	DismissPanicFunc(func(PanicEvent)) EventListener

	OnSuccess(ge.Action[SuccessEvent]) EventListener
	OnTimeout(ge.Action[TimeoutEvent]) EventListener
	OnPanic(ge.Action[PanicEvent]) EventListener
	DismissSuccess(ge.Action[SuccessEvent]) EventListener
	DismissTimeout(ge.Action[TimeoutEvent]) EventListener
	DismissPanic(ge.Action[PanicEvent]) EventListener
}

func newEventListener() *eventListener {
	return &eventListener{
		onSuccess: make([]ge.Action[SuccessEvent], 0),
		onTimeout: make([]ge.Action[TimeoutEvent], 0),
		onPanic:   make([]ge.Action[PanicEvent], 0),
	}
}

type eventListener struct {
	sync.RWMutex
	onSuccess []ge.Action[SuccessEvent]
	onTimeout []ge.Action[TimeoutEvent]
	onPanic   []ge.Action[PanicEvent]
}

func (listener *eventListener) OnSuccessFunc(consumer func(SuccessEvent)) EventListener {
	return listener.OnSuccess(ge.ActionFunc[SuccessEvent](consumer))
}

func (listener *eventListener) OnTimeoutFunc(consumer func(TimeoutEvent)) EventListener {
	return listener.OnTimeout(ge.ActionFunc[TimeoutEvent](consumer))
}

func (listener *eventListener) OnPanicFunc(consumer func(PanicEvent)) EventListener {
	return listener.OnPanic(ge.ActionFunc[PanicEvent](consumer))
}

func (listener *eventListener) DismissSuccessFunc(consumer func(SuccessEvent)) EventListener {
	return listener.DismissSuccess(ge.ActionFunc[SuccessEvent](consumer))
}

func (listener *eventListener) DismissTimeoutFunc(consumer func(TimeoutEvent)) EventListener {
	return listener.DismissTimeout(ge.ActionFunc[TimeoutEvent](consumer))
}

func (listener *eventListener) DismissPanicFunc(consumer func(PanicEvent)) EventListener {
	return listener.DismissPanic(ge.ActionFunc[PanicEvent](consumer))
}

func (listener *eventListener) OnSuccess(action ge.Action[SuccessEvent]) EventListener {
	listener.Lock()
	defer listener.Unlock()
	listener.onSuccess = ge.AppendElementUnique(listener.onSuccess, action)
	return listener
}

func (listener *eventListener) OnTimeout(action ge.Action[TimeoutEvent]) EventListener {
	listener.Lock()
	defer listener.Unlock()
	listener.onTimeout = ge.AppendElementUnique(listener.onTimeout, action)
	return listener
}

func (listener *eventListener) OnPanic(action ge.Action[PanicEvent]) EventListener {
	listener.Lock()
	defer listener.Unlock()
	listener.onPanic = ge.AppendElementUnique(listener.onPanic, action)
	return listener
}

func (listener *eventListener) DismissSuccess(action ge.Action[SuccessEvent]) EventListener {
	listener.Lock()
	defer listener.Unlock()
	listener.onSuccess = ge.RemoveElementByValue(listener.onSuccess, action)
	return listener
}

func (listener *eventListener) DismissTimeout(action ge.Action[TimeoutEvent]) EventListener {
	listener.Lock()
	defer listener.Unlock()
	listener.onTimeout = ge.RemoveElementByValue(listener.onTimeout, action)
	return listener
}

func (listener *eventListener) DismissPanic(action ge.Action[PanicEvent]) EventListener {
	listener.Lock()
	defer listener.Unlock()
	listener.onPanic = ge.RemoveElementByValue(listener.onPanic, action)
	return listener
}

func (listener *eventListener) consumeEvent(event Event) {
	go func() {
		listener.RLock()
		defer listener.RUnlock()
		switch e := event.(type) {
		case *successEvent:
			ge.ForEach(listener.onSuccess, SuccessEvent(e))
		case *timeoutEvent:
			ge.ForEach(listener.onTimeout, TimeoutEvent(e))
		case *panicEvent:
			ge.ForEach(listener.onPanic, PanicEvent(e))
		}
	}()
}
