package circuitbreaker

import (
	"fmt"
	"github.com/CharLemAznable/gofn/common"
	"sync/atomic"
	"time"
)

type CircuitBreaker interface {
	Name() string
	EventListener() EventListener
	TransitionToDisabled() error
	TransitionToForcedOpen() error
	TransitionToClosedState() error
	TransitionToOpenState() error
	TransitionToHalfOpenState() error
	Metrics() Metrics

	config() *Config
	execute(func() (any, error)) (any, error)
	acquirePermission() error
	publishThresholdsExceededEvent(metricsResult, Metrics)
}

func NewCircuitBreaker(name string, configs ...ConfigBuilder) CircuitBreaker {
	config := defaultConfig()
	for _, cfg := range configs {
		cfg(config)
	}
	machine := &stateMachine{
		name:          name,
		conf:          config,
		eventListener: newEventListener(),
	}
	machine.state.Store(closed(machine))
	return machine
}

type stateMachine struct {
	name          string
	conf          *Config
	eventListener EventListener
	state         atomic.Pointer[state]
}

func (machine *stateMachine) Name() string {
	return machine.name
}

func (machine *stateMachine) EventListener() EventListener {
	return machine.eventListener
}

func (machine *stateMachine) TransitionToDisabled() error {
	return machine.stateTransition(Disabled, func(_ *state) *state {
		return disabled(machine)
	})
}

func (machine *stateMachine) TransitionToForcedOpen() error {
	return machine.stateTransition(ForcedOpen, func(current *state) *state {
		return forcedOpen(current.attempts, machine)
	})
}

func (machine *stateMachine) TransitionToClosedState() error {
	return machine.stateTransition(Closed, func(_ *state) *state {
		return closed(machine)
	})
}

func (machine *stateMachine) TransitionToOpenState() error {
	return machine.stateTransition(Open, func(current *state) *state {
		return open(current.attempts+1, current.metrics, machine)
	})
}

func (machine *stateMachine) TransitionToHalfOpenState() error {
	return machine.stateTransition(HalfOpen, func(current *state) *state {
		return halfOpen(current.attempts, machine)
	})
}

func (machine *stateMachine) stateTransition(newStateName stateName, generator func(*state) *state) error {
	var previous *state
	previous, err := getAndUpdatePointer(&machine.state, func(currentState *state) (*state, error) {
		if err := checkStateTransition(machine.name, currentState.name, newStateName); err != nil {
			return nil, err
		}
		if currentState.preTransitionHook != nil {
			go currentState.preTransitionHook()
		}
		return generator(currentState), nil
	})
	if err == nil && previous.name != newStateName {
		machine.publishEvent(newStateTransitionEvent(machine.name, previous.name, newStateName))
	}
	return err
}

func (machine *stateMachine) Metrics() Metrics {
	return machine.loadState().metrics
}

func (machine *stateMachine) config() *Config {
	return machine.conf
}

func (machine *stateMachine) execute(fn func() (any, error)) (any, error) {
	if err := machine.acquirePermission(); err != nil {
		machine.publishEvent(newNotPermittedEvent(machine.name))
		return nil, err
	}
	start := time.Now()
	finished := make(chan *channelValue)
	panicked := make(common.Panicked)
	go func() {
		defer panicked.Recover()
		ret, err := fn()
		finished <- &channelValue{ret, err}
	}()
	select {
	case result := <-finished:
		machine.onResult(start, result.ret, result.err)
		return result.ret, result.err
	case err := <-panicked.Caught():
		machine.onResult(start, nil, common.WrapPanic(err))
		panic(err)
	}
}

func (machine *stateMachine) acquirePermission() error {
	if fn := machine.loadState().acquirePermission; fn != nil {
		return fn()
	}
	return nil
}

func (machine *stateMachine) onResult(start time.Time, ret any, err error) {
	duration := time.Now().Sub(start)
	if machine.conf.recordResultPredicateFn(ret, err) {
		machine.publishEvent(newErrorEvent(machine.name, duration, ret, err))
		if fn := machine.loadState().onError; fn != nil {
			fn(duration)
		}
	} else {
		machine.publishEvent(newSuccessEvent(machine.name, duration))
		if fn := machine.loadState().onSuccess; fn != nil {
			fn(duration)
		}
	}
}

func (machine *stateMachine) publishThresholdsExceededEvent(result metricsResult, metrics Metrics) {
	if failureRateExceededThreshold(result) {
		machine.publishEvent(newFailureRateExceededEvent(machine.name, metrics.FailureRate()))
	}
	if slowCallRateExceededThreshold(result) {
		machine.publishEvent(newSlowCallRateExceededEvent(machine.name, metrics.SlowCallRate()))
	}
}

func (machine *stateMachine) loadState() *state {
	return machine.state.Load()
}

func (machine *stateMachine) publishEvent(event Event) {
	if event.EventType() == StateTransition ||
		machine.loadState().allowPublish {
		machine.eventListener.consumeEvent(event)
	}
}

type channelValue struct {
	ret any
	err error
}

type NotPermittedError struct {
	name      string
	stateName stateName
}

func (e *NotPermittedError) Error() string {
	return fmt.Sprintf("CircuitBreaker '%s' is %s and does not permit further calls", e.name, e.stateName)
}
