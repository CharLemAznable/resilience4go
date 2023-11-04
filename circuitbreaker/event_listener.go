package circuitbreaker

type EventConsumer func(Event)

type EventListener interface {
	OnSuccess(EventConsumer) EventListener
	OnError(EventConsumer) EventListener
	OnNotPermitted(EventConsumer) EventListener
	OnStateTransition(EventConsumer) EventListener
	OnFailureRateExceeded(EventConsumer) EventListener
	OnSlowCallRateExceeded(EventConsumer) EventListener
	consumeEvent(Event)
}

func newEventListener() EventListener {
	return &eventListener{
		onSuccess:              make([]EventConsumer, 0),
		onError:                make([]EventConsumer, 0),
		onNotPermitted:         make([]EventConsumer, 0),
		onStateTransition:      make([]EventConsumer, 0),
		onFailureRateExceeded:  make([]EventConsumer, 0),
		onSlowCallRateExceeded: make([]EventConsumer, 0),
	}
}

type eventListener struct {
	onSuccess              []EventConsumer
	onError                []EventConsumer
	onNotPermitted         []EventConsumer
	onStateTransition      []EventConsumer
	onFailureRateExceeded  []EventConsumer
	onSlowCallRateExceeded []EventConsumer
}

func (listener *eventListener) OnSuccess(consumer EventConsumer) EventListener {
	listener.onSuccess = append(listener.onSuccess, consumer)
	return listener
}

func (listener *eventListener) OnError(consumer EventConsumer) EventListener {
	listener.onError = append(listener.onError, consumer)
	return listener
}

func (listener *eventListener) OnNotPermitted(consumer EventConsumer) EventListener {
	listener.onNotPermitted = append(listener.onNotPermitted, consumer)
	return listener
}

func (listener *eventListener) OnStateTransition(consumer EventConsumer) EventListener {
	listener.onStateTransition = append(listener.onStateTransition, consumer)
	return listener
}

func (listener *eventListener) OnFailureRateExceeded(consumer EventConsumer) EventListener {
	listener.onFailureRateExceeded = append(listener.onFailureRateExceeded, consumer)
	return listener
}

func (listener *eventListener) OnSlowCallRateExceeded(consumer EventConsumer) EventListener {
	listener.onSlowCallRateExceeded = append(listener.onSlowCallRateExceeded, consumer)
	return listener
}

func (listener *eventListener) consumeEvent(event Event) {
	var consumers []EventConsumer
	switch event.EventType() {
	case Success:
		consumers = listener.onSuccess
	case Error:
		consumers = listener.onError
	case NotPermitted:
		consumers = listener.onNotPermitted
	case StateTransition:
		consumers = listener.onStateTransition
	case FailureRateExceeded:
		consumers = listener.onFailureRateExceeded
	case SlowCallRateExceeded:
		consumers = listener.onSlowCallRateExceeded
	}
	for _, consumer := range consumers {
		go consumer(event)
	}
}
