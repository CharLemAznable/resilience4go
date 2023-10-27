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

func (processor *eventListener) OnPermitted(consumer EventConsumer) EventListener {
	processor.onPermitted = append(processor.onPermitted, consumer)
	return processor
}

func (processor *eventListener) OnRejected(consumer EventConsumer) EventListener {
	processor.onRejected = append(processor.onRejected, consumer)
	return processor
}

func (processor *eventListener) OnFinished(consumer EventConsumer) EventListener {
	processor.onFinished = append(processor.onFinished, consumer)
	return processor
}

func (processor *eventListener) consumeEvent(event Event) {
	var consumers []EventConsumer
	switch event.EventType() {
	case PERMITTED:
		consumers = processor.onPermitted
	case REJECTED:
		consumers = processor.onRejected
	case FINISHED:
		consumers = processor.onFinished
	}
	for _, consumer := range consumers {
		go consumer(event)
	}
}
