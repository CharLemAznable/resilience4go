package ratelimiter

import (
	"fmt"
	"github.com/CharLemAznable/resilience4go/utils"
	"sync/atomic"
	"time"
)

type RateLimiter interface {
	Name() string
	Metrics() Metrics
	EventListener() EventListener

	acquirePermission() error
}

func NewRateLimiter(name string, configs ...ConfigBuilder) RateLimiter {
	config := defaultConfig()
	for _, cfg := range configs {
		cfg(config)
	}
	limiter := &atomicRateLimiter{
		name:          name,
		config:        config,
		nanoTimeStart: time.Now().UnixNano(),
		eventListener: newEventListener(),
	}
	limiter.state.Store(&state{
		activeCycle:       0,
		activePermissions: config.limitForPeriod,
		nanosToWait:       0,
	})
	limiter.metrics = newMetrics(
		func() int64 {
			return limiter.waitingThreads.Load()
		},
		func() int64 {
			currentState := limiter.state.Load()
			estimatedState := limiter.calculateNextState(-1, currentState)
			return estimatedState.activePermissions
		})
	return limiter
}

type atomicRateLimiter struct {
	name           string
	config         *Config
	nanoTimeStart  int64
	state          atomic.Pointer[state]
	waitingThreads atomic.Int64
	metrics        Metrics
	eventListener  EventListener
}

func (limiter *atomicRateLimiter) Name() string {
	return limiter.name
}

func (limiter *atomicRateLimiter) Metrics() Metrics {
	return limiter.metrics
}

func (limiter *atomicRateLimiter) EventListener() EventListener {
	return limiter.eventListener
}

func (limiter *atomicRateLimiter) acquirePermission() error {
	timeoutInNanos := limiter.config.timeoutDuration.Nanoseconds()
	modifiedState := limiter.updateStateWithBackOff(timeoutInNanos)
	if limiter.waitForPermissionIfNecessary(timeoutInNanos, modifiedState.nanosToWait) {
		limiter.eventListener.consumeEvent(newSuccessEvent(limiter.name))
		return nil
	}
	limiter.eventListener.consumeEvent(newFailureEvent(limiter.name))
	return &NotPermittedError{name: limiter.name}
}

func (limiter *atomicRateLimiter) updateStateWithBackOff(timeoutInNanos int64) *state {
	for {
		prev := limiter.state.Load()
		next := limiter.calculateNextState(timeoutInNanos, prev)
		if limiter.compareAndSet(prev, next) {
			return next
		}
	}
}

func (limiter *atomicRateLimiter) compareAndSet(current, next *state) bool {
	if limiter.state.CompareAndSwap(current, next) {
		return true
	}
	time.Sleep(1)
	return false
}

func (limiter *atomicRateLimiter) calculateNextState(timeoutInNanos int64, activeState *state) *state {
	cyclePeriodInNanos := limiter.config.limitRefreshPeriod.Nanoseconds()
	permissionsPerCycle := limiter.config.limitForPeriod
	currentNanos := limiter.currentNanoTime()
	currentCycle := currentNanos / cyclePeriodInNanos
	nextCycle := activeState.activeCycle
	nextPermissions := activeState.activePermissions
	if nextCycle != currentCycle {
		elapsedCycles := currentCycle - nextCycle
		accumulatedPermissions := elapsedCycles * permissionsPerCycle
		nextCycle = currentCycle
		nextPermissions = utils.Min(nextPermissions+accumulatedPermissions, permissionsPerCycle)
	}
	nextNanosToWait := nanosToWaitForPermission(cyclePeriodInNanos,
		permissionsPerCycle, nextPermissions, currentNanos, currentCycle)
	return reservePermissions(timeoutInNanos, nextCycle, nextPermissions, nextNanosToWait)
}

func (limiter *atomicRateLimiter) waitForPermissionIfNecessary(
	timeoutInNanos, nanosToWait int64) bool {
	if nanosToWait <= 0 {
		return true
	}
	if timeoutInNanos >= nanosToWait {
		limiter.waitForPermission(nanosToWait)
		return true
	}
	limiter.waitForPermission(timeoutInNanos)
	return false
}

func (limiter *atomicRateLimiter) waitForPermission(nanosToWait int64) {
	limiter.waitingThreads.Add(1)
	deadline := limiter.currentNanoTime() + nanosToWait
	for limiter.currentNanoTime() < deadline {
		time.Sleep(time.Duration(deadline - limiter.currentNanoTime()))
	}
	limiter.waitingThreads.Add(-1)
}

func (limiter *atomicRateLimiter) currentNanoTime() int64 {
	return time.Now().UnixNano() - limiter.nanoTimeStart
}

func nanosToWaitForPermission(
	cyclePeriodInNanos, permissionsPerCycle,
	availablePermissions, currentNanos, currentCycle int64) int64 {
	if availablePermissions >= 1 {
		return 0
	}
	nextCycleTimeInNanos := (currentCycle + 1) * cyclePeriodInNanos
	nanosToNextCycle := nextCycleTimeInNanos - currentNanos
	permissionsAtTheStartOfNextCycle := availablePermissions + permissionsPerCycle
	fullCyclesToWait := divCeil(-(permissionsAtTheStartOfNextCycle - 1), permissionsPerCycle)
	return (fullCyclesToWait * cyclePeriodInNanos) + nanosToNextCycle
}

func divCeil(x, y int64) int64 {
	return (x + y - 1) / y
}

func reservePermissions(
	timeoutInNanos, cycle, permissions, nanosToWait int64) *state {
	permissionsWithReservation := permissions
	if timeoutInNanos >= nanosToWait {
		permissionsWithReservation -= 1
	}
	return &state{
		activeCycle:       cycle,
		activePermissions: permissionsWithReservation,
		nanosToWait:       nanosToWait}
}

type state struct {
	activeCycle       int64
	activePermissions int64
	nanosToWait       int64
}

type NotPermittedError struct {
	name string
}

func (e *NotPermittedError) Error() string {
	return fmt.Sprintf("RateLimiter '%s' does not permit further calls", e.name)
}
