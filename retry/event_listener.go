package retry

import (
	"github.com/CharLemAznable/resilience4go/utils"
	"sync"
)

type SuccessEventConsumer func(SuccessEvent)

//goland:noinspection GoNameStartsWithPackageName
type RetryEventConsumer func(RetryEvent)

type ErrorEventConsumer func(ErrorEvent)

type EventListener interface {
	OnSuccess(SuccessEventConsumer) EventListener
	OnRetry(RetryEventConsumer) EventListener
	OnError(ErrorEventConsumer) EventListener
	Dismiss(any) EventListener
	consumeEvent(Event)
}

func newEventListener() EventListener {
	return &eventListener{
		onSuccess:       make([]SuccessEventConsumer, 0),
		onSuccessSlices: utils.NewSlicesWithPointer[SuccessEventConsumer](),
		onRetry:         make([]RetryEventConsumer, 0),
		onRetrySlices:   utils.NewSlicesWithPointer[RetryEventConsumer](),
		onError:         make([]ErrorEventConsumer, 0),
		onErrorSlices:   utils.NewSlicesWithPointer[ErrorEventConsumer](),
	}
}

type eventListener struct {
	sync.RWMutex
	onSuccess       []SuccessEventConsumer
	onSuccessSlices utils.Slices[SuccessEventConsumer]
	onRetry         []RetryEventConsumer
	onRetrySlices   utils.Slices[RetryEventConsumer]
	onError         []ErrorEventConsumer
	onErrorSlices   utils.Slices[ErrorEventConsumer]
}

func (listener *eventListener) OnSuccess(consumer SuccessEventConsumer) EventListener {
	listener.Lock()
	defer listener.Unlock()
	listener.onSuccess = listener.onSuccessSlices.AppendElementUnique(listener.onSuccess, consumer)
	return listener
}

func (listener *eventListener) OnRetry(consumer RetryEventConsumer) EventListener {
	listener.Lock()
	defer listener.Unlock()
	listener.onRetry = listener.onRetrySlices.AppendElementUnique(listener.onRetry, consumer)
	return listener
}

func (listener *eventListener) OnError(consumer ErrorEventConsumer) EventListener {
	listener.Lock()
	defer listener.Unlock()
	listener.onError = listener.onErrorSlices.AppendElementUnique(listener.onError, consumer)
	return listener
}

func (listener *eventListener) Dismiss(consumer any) EventListener {
	listener.Lock()
	defer listener.Unlock()
	if c, ok := consumer.(func(SuccessEvent)); ok {
		listener.onSuccess = listener.onSuccessSlices.RemoveElementByValue(listener.onSuccess, c)
	}
	if c, ok := consumer.(func(RetryEvent)); ok {
		listener.onRetry = listener.onRetrySlices.RemoveElementByValue(listener.onRetry, c)
	}
	if c, ok := consumer.(func(ErrorEvent)); ok {
		listener.onError = listener.onErrorSlices.RemoveElementByValue(listener.onError, c)
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
		case *retryEvent:
			for _, consumer := range listener.onRetry {
				go consumer(e)
			}
		case *errorEvent:
			for _, consumer := range listener.onError {
				go consumer(e)
			}
		}
	}()
}
