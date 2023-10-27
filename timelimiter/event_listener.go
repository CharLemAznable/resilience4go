package timelimiter

type EventConsumer func(Event)

type EventListener interface {
	OnSuccess(EventConsumer) EventListener
	OnTimeout(EventConsumer) EventListener
	OnFailure(EventConsumer) EventListener
	consumeEvent(Event)
}

func newEventListener() EventListener {
	return &eventListener{
		onSuccess: make([]EventConsumer, 0),
		onTimeout: make([]EventConsumer, 0),
		onFailure: make([]EventConsumer, 0),
	}
}

type eventListener struct {
	onSuccess []EventConsumer
	onTimeout []EventConsumer
	onFailure []EventConsumer
}

func (listener *eventListener) OnSuccess(consumer EventConsumer) EventListener {
	listener.onSuccess = append(listener.onSuccess, consumer)
	return listener
}

func (listener *eventListener) OnTimeout(consumer EventConsumer) EventListener {
	listener.onTimeout = append(listener.onTimeout, consumer)
	return listener
}

func (listener *eventListener) OnFailure(consumer EventConsumer) EventListener {
	listener.onFailure = append(listener.onFailure, consumer)
	return listener
}

func (listener *eventListener) consumeEvent(event Event) {
	var consumers []EventConsumer
	switch event.EventType() {
	case SUCCESS:
		consumers = listener.onSuccess
	case TIMEOUT:
		consumers = listener.onTimeout
	case FAILURE:
		consumers = listener.onFailure
	}
	for _, consumer := range consumers {
		go consumer(event)
	}
}
