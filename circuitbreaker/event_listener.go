package circuitbreaker

import (
	"github.com/CharLemAznable/resilience4go/utils"
	"sync"
)

type SuccessEventConsumer func(SuccessEvent)
type ErrorEventConsumer func(ErrorEvent)
type NotPermittedEventConsumer func(NotPermittedEvent)
type StateTransitionEventConsumer func(StateTransitionEvent)
type FailureRateExceededEventConsumer func(FailureRateExceededEvent)
type SlowCallRateExceededEventConsumer func(SlowCallRateExceededEvent)

type EventListener interface {
	OnSuccess(SuccessEventConsumer) EventListener
	OnError(ErrorEventConsumer) EventListener
	OnNotPermitted(NotPermittedEventConsumer) EventListener
	OnStateTransition(StateTransitionEventConsumer) EventListener
	OnFailureRateExceeded(FailureRateExceededEventConsumer) EventListener
	OnSlowCallRateExceeded(SlowCallRateExceededEventConsumer) EventListener
	Dismiss(any) EventListener
	consumeEvent(Event)
}

func newEventListener() EventListener {
	return &eventListener{
		onSuccess:                    make([]SuccessEventConsumer, 0),
		onSuccessSlices:              utils.NewSlicesWithPointer[SuccessEventConsumer](),
		onError:                      make([]ErrorEventConsumer, 0),
		onErrorSlices:                utils.NewSlicesWithPointer[ErrorEventConsumer](),
		onNotPermitted:               make([]NotPermittedEventConsumer, 0),
		onNotPermittedSlices:         utils.NewSlicesWithPointer[NotPermittedEventConsumer](),
		onStateTransition:            make([]StateTransitionEventConsumer, 0),
		onStateTransitionSlices:      utils.NewSlicesWithPointer[StateTransitionEventConsumer](),
		onFailureRateExceeded:        make([]FailureRateExceededEventConsumer, 0),
		onFailureRateExceededSlices:  utils.NewSlicesWithPointer[FailureRateExceededEventConsumer](),
		onSlowCallRateExceeded:       make([]SlowCallRateExceededEventConsumer, 0),
		onSlowCallRateExceededSlices: utils.NewSlicesWithPointer[SlowCallRateExceededEventConsumer](),
	}
}

type eventListener struct {
	sync.RWMutex
	onSuccess                    []SuccessEventConsumer
	onSuccessSlices              utils.Slices[SuccessEventConsumer]
	onError                      []ErrorEventConsumer
	onErrorSlices                utils.Slices[ErrorEventConsumer]
	onNotPermitted               []NotPermittedEventConsumer
	onNotPermittedSlices         utils.Slices[NotPermittedEventConsumer]
	onStateTransition            []StateTransitionEventConsumer
	onStateTransitionSlices      utils.Slices[StateTransitionEventConsumer]
	onFailureRateExceeded        []FailureRateExceededEventConsumer
	onFailureRateExceededSlices  utils.Slices[FailureRateExceededEventConsumer]
	onSlowCallRateExceeded       []SlowCallRateExceededEventConsumer
	onSlowCallRateExceededSlices utils.Slices[SlowCallRateExceededEventConsumer]
}

func (listener *eventListener) OnSuccess(consumer SuccessEventConsumer) EventListener {
	listener.Lock()
	defer listener.Unlock()
	listener.onSuccess = listener.onSuccessSlices.AppendElementUnique(listener.onSuccess, consumer)
	return listener
}

func (listener *eventListener) OnError(consumer ErrorEventConsumer) EventListener {
	listener.Lock()
	defer listener.Unlock()
	listener.onError = listener.onErrorSlices.AppendElementUnique(listener.onError, consumer)
	return listener
}

func (listener *eventListener) OnNotPermitted(consumer NotPermittedEventConsumer) EventListener {
	listener.Lock()
	defer listener.Unlock()
	listener.onNotPermitted = listener.onNotPermittedSlices.AppendElementUnique(listener.onNotPermitted, consumer)
	return listener
}

func (listener *eventListener) OnStateTransition(consumer StateTransitionEventConsumer) EventListener {
	listener.Lock()
	defer listener.Unlock()
	listener.onStateTransition = listener.onStateTransitionSlices.AppendElementUnique(listener.onStateTransition, consumer)
	return listener
}

func (listener *eventListener) OnFailureRateExceeded(consumer FailureRateExceededEventConsumer) EventListener {
	listener.Lock()
	defer listener.Unlock()
	listener.onFailureRateExceeded = listener.onFailureRateExceededSlices.AppendElementUnique(listener.onFailureRateExceeded, consumer)
	return listener
}

func (listener *eventListener) OnSlowCallRateExceeded(consumer SlowCallRateExceededEventConsumer) EventListener {
	listener.Lock()
	defer listener.Unlock()
	listener.onSlowCallRateExceeded = listener.onSlowCallRateExceededSlices.AppendElementUnique(listener.onSlowCallRateExceeded, consumer)
	return listener
}

func (listener *eventListener) Dismiss(consumer any) EventListener {
	listener.Lock()
	defer listener.Unlock()
	if c, ok := consumer.(func(SuccessEvent)); ok {
		listener.onSuccess = listener.onSuccessSlices.RemoveElementByValue(listener.onSuccess, c)
	}
	if c, ok := consumer.(func(ErrorEvent)); ok {
		listener.onError = listener.onErrorSlices.RemoveElementByValue(listener.onError, c)
	}
	if c, ok := consumer.(func(NotPermittedEvent)); ok {
		listener.onNotPermitted = listener.onNotPermittedSlices.RemoveElementByValue(listener.onNotPermitted, c)
	}
	if c, ok := consumer.(func(StateTransitionEvent)); ok {
		listener.onStateTransition = listener.onStateTransitionSlices.RemoveElementByValue(listener.onStateTransition, c)
	}
	if c, ok := consumer.(func(FailureRateExceededEvent)); ok {
		listener.onFailureRateExceeded = listener.onFailureRateExceededSlices.RemoveElementByValue(listener.onFailureRateExceeded, c)
	}
	if c, ok := consumer.(func(SlowCallRateExceededEvent)); ok {
		listener.onSlowCallRateExceeded = listener.onSlowCallRateExceededSlices.RemoveElementByValue(listener.onSlowCallRateExceeded, c)
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
		case *errorEvent:
			for _, consumer := range listener.onError {
				go consumer(e)
			}
		case *notPermittedEvent:
			for _, consumer := range listener.onNotPermitted {
				go consumer(e)
			}
		case *stateTransitionEvent:
			for _, consumer := range listener.onStateTransition {
				go consumer(e)
			}
		case *failureRateExceededEvent:
			for _, consumer := range listener.onFailureRateExceeded {
				go consumer(e)
			}
		case *slowCallRateExceededEvent:
			for _, consumer := range listener.onSlowCallRateExceeded {
				go consumer(e)
			}
		}
	}()
}
