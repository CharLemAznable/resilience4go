package bulkhead

type EventConsumer func(Event)

type EventListener interface {
	OnPermitted(EventConsumer) EventListener
	OnRejected(EventConsumer) EventListener
	OnFinished(EventConsumer) EventListener
	consumeEvent(Event)
}

func newEventListener() EventListener {
	return &eventListener{
		onPermitted: make([]EventConsumer, 0),
		onRejected:  make([]EventConsumer, 0),
		onFinished:  make([]EventConsumer, 0),
	}
}

type eventListener struct {
	onPermitted []EventConsumer
	onRejected  []EventConsumer
	onFinished  []EventConsumer
}

func (listener *eventListener) OnPermitted(consumer EventConsumer) EventListener {
	listener.onPermitted = append(listener.onPermitted, consumer)
	return listener
}

func (listener *eventListener) OnRejected(consumer EventConsumer) EventListener {
	listener.onRejected = append(listener.onRejected, consumer)
	return listener
}

func (listener *eventListener) OnFinished(consumer EventConsumer) EventListener {
	listener.onFinished = append(listener.onFinished, consumer)
	return listener
}

func (listener *eventListener) consumeEvent(event Event) {
	var consumers []EventConsumer
	switch event.EventType() {
	case PERMITTED:
		consumers = listener.onPermitted
	case REJECTED:
		consumers = listener.onRejected
	case FINISHED:
		consumers = listener.onFinished
	}
	for _, consumer := range consumers {
		go consumer(event)
	}
}
