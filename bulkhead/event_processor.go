package bulkhead

type EventConsumer func(Event)

type EventProcessor interface {
	OnPermitted(EventConsumer) EventProcessor
	OnRejected(EventConsumer) EventProcessor
	OnFinished(EventConsumer) EventProcessor
	consumeEvent(Event)
}

func newEventProcessor() EventProcessor {
	return &eventProcessor{
		onPermitted: make([]EventConsumer, 0),
		onRejected:  make([]EventConsumer, 0),
		onFinished:  make([]EventConsumer, 0),
	}
}

type eventProcessor struct {
	onPermitted []EventConsumer
	onRejected  []EventConsumer
	onFinished  []EventConsumer
}

func (processor *eventProcessor) OnPermitted(consumer EventConsumer) EventProcessor {
	processor.onPermitted = append(processor.onPermitted, consumer)
	return processor
}

func (processor *eventProcessor) OnRejected(consumer EventConsumer) EventProcessor {
	processor.onRejected = append(processor.onRejected, consumer)
	return processor
}

func (processor *eventProcessor) OnFinished(consumer EventConsumer) EventProcessor {
	processor.onFinished = append(processor.onFinished, consumer)
	return processor
}

func (processor *eventProcessor) consumeEvent(event Event) {
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
