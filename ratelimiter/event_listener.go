package ratelimiter

import (
	"github.com/CharLemAznable/resilience4go/utils"
	"sync"
)

type SuccessEventConsumer func(SuccessEvent)
type FailureEventConsumer func(FailureEvent)

type EventListener interface {
	OnSuccess(SuccessEventConsumer) EventListener
	OnFailure(FailureEventConsumer) EventListener
	Dismiss(any) EventListener
	consumeEvent(Event)
}

func newEventListener() EventListener {
	return &eventListener{
		onSuccess:       make([]SuccessEventConsumer, 0),
		onSuccessSlices: utils.NewSlicesWithPointer[SuccessEventConsumer](),
		onFailure:       make([]FailureEventConsumer, 0),
		onFailureSlices: utils.NewSlicesWithPointer[FailureEventConsumer](),
	}
}

type eventListener struct {
	sync.RWMutex
	onSuccess       []SuccessEventConsumer
	onSuccessSlices utils.Slices[SuccessEventConsumer]
	onFailure       []FailureEventConsumer
	onFailureSlices utils.Slices[FailureEventConsumer]
}

func (listener *eventListener) OnSuccess(consumer SuccessEventConsumer) EventListener {
	listener.Lock()
	defer listener.Unlock()
	listener.onSuccess = listener.onSuccessSlices.AppendElementUnique(listener.onSuccess, consumer)
	return listener
}

func (listener *eventListener) OnFailure(consumer FailureEventConsumer) EventListener {
	listener.Lock()
	defer listener.Unlock()
	listener.onFailure = listener.onFailureSlices.AppendElementUnique(listener.onFailure, consumer)
	return listener
}

func (listener *eventListener) Dismiss(consumer any) EventListener {
	listener.Lock()
	defer listener.Unlock()
	if c, ok := consumer.(func(SuccessEvent)); ok {
		listener.onSuccess = listener.onSuccessSlices.RemoveElementByValue(listener.onSuccess, c)
	}
	if c, ok := consumer.(func(FailureEvent)); ok {
		listener.onFailure = listener.onFailureSlices.RemoveElementByValue(listener.onFailure, c)
	}
	return listener
}

func (listener *eventListener) consumeEvent(event Event) {
	go func() {
		listener.RLock()
		defer listener.RUnlock()
		switch e := event.(type) {
		case *successEvent:
			for _, consumer := range listener.onSuccess {
				go consumer(e)
			}
		case *failureEvent:
			for _, consumer := range listener.onFailure {
				go consumer(e)
			}
		}
	}()
}
