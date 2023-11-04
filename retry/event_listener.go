package retry

type EventConsumer func(Event)

type EventListener interface {
	OnSuccess(EventConsumer) EventListener
	OnError(EventConsumer) EventListener
	OnRetry(EventConsumer) EventListener
	consumeEvent(Event)
}

func newEventListener() EventListener {
	return &eventListener{
		onSuccess: make([]EventConsumer, 0),
		onError:   make([]EventConsumer, 0),
		onRetry:   make([]EventConsumer, 0),
	}
}

type eventListener struct {
	onSuccess []EventConsumer
	onError   []EventConsumer
	onRetry   []EventConsumer
}

func (listener *eventListener) OnSuccess(consumer EventConsumer) EventListener {
	listener.onSuccess = append(listener.onSuccess, consumer)
	return listener
}

func (listener *eventListener) OnError(consumer EventConsumer) EventListener {
	listener.onError = append(listener.onError, consumer)
	return listener
}

func (listener *eventListener) OnRetry(consumer EventConsumer) EventListener {
	listener.onRetry = append(listener.onRetry, consumer)
	return listener
}

func (listener *eventListener) consumeEvent(event Event) {
	var consumers []EventConsumer
	switch event.EventType() {
	case SUCCESS:
		consumers = listener.onSuccess
	case ERROR:
		consumers = listener.onError
	case RETRY:
		consumers = listener.onRetry
	}
	for _, consumer := range consumers {
		go consumer(event)
	}
}
