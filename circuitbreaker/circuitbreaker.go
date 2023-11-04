package circuitbreaker

import (
	"fmt"
	"github.com/CharLemAznable/gofn/common"
	"sync/atomic"
	"time"
)

type CircuitBreaker interface {
	Name() string
	TransitionToDisabled() error
	TransitionToForcedOpen() error
	TransitionToClosedState() error
	TransitionToOpenState() error
	TransitionToHalfOpenState() error
	Metrics() Metrics

	config() *Config
	execute(fn func() (any, error)) (any, error)
	acquirePermission() error
}

func NewCircuitBreaker(name string, configs ...ConfigBuilder) CircuitBreaker {
	config := defaultConfig()
	for _, cfg := range configs {
		cfg(config)
	}
	machine := &stateMachine{
		name: name,
		conf: config,
	}
	machine.state.Store(closed(machine))
	return machine
}

type stateMachine struct {
	name  string
	conf  *Config
	state atomic.Pointer[state]
}

func (machine *stateMachine) Name() string {
	return machine.name
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

func (machine *stateMachine) Metrics() Metrics {
	return machine.loadState().metrics
}

func (machine *stateMachine) config() *Config {
	return machine.conf
}

func (machine *stateMachine) execute(fn func() (any, error)) (any, error) {
	if err := machine.acquirePermission(); err != nil {
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
		if fn := machine.loadState().onError; fn != nil {
			fn(duration)
		}
	} else {
		if fn := machine.loadState().onSuccess; fn != nil {
			fn(duration)
		}
	}
}

func (machine *stateMachine) loadState() *state {
	return machine.state.Load()
}

func (machine *stateMachine) stateTransition(newStateName stateName, generator func(*state) *state) error {
	_, err := getAndUpdatePointer(&machine.state, func(currentState *state) (*state, error) {
		if _, err := newStateTransition(machine.name, currentState.name, newStateName); err != nil {
			return nil, err
		}
		if currentState.preTransitionHook != nil {
			go currentState.preTransitionHook()
		}
		return generator(currentState), nil
	})
	return err
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
