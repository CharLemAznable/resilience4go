package ratelimiter

type EventConsumer func(Event)

type EventListener interface {
	OnSuccess(EventConsumer) EventListener
	OnFailure(EventConsumer) EventListener
	consumeEvent(Event)
}

func newEventListener() EventListener {
	return &eventListener{
		onSuccess: make([]EventConsumer, 0),
		onFailure: make([]EventConsumer, 0),
	}
}

type eventListener struct {
	onSuccess []EventConsumer
	onFailure []EventConsumer
}

func (listener *eventListener) OnSuccess(consumer EventConsumer) EventListener {
	listener.onSuccess = append(listener.onSuccess, consumer)
	return listener
}

func (listener *eventListener) OnFailure(consumer EventConsumer) EventListener {
	listener.onFailure = append(listener.onFailure, consumer)
	return listener
}

func (listener *eventListener) consumeEvent(event Event) {
	var consumers []EventConsumer
	switch event.EventType() {
	case SUCCESSFUL:
		consumers = listener.onSuccess
	case FAILED:
		consumers = listener.onFailure
	}
	for _, consumer := range consumers {
		go consumer(event)
	}
}
