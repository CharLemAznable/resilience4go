package circuitbreaker_test

import (
	"github.com/CharLemAznable/resilience4go/circuitbreaker"
	"testing"
)

func TestDecorateCover(t *testing.T) {
	breaker := circuitbreaker.NewCircuitBreaker("test")
	circuitbreaker.DecorateRun(breaker, func() {})
	circuitbreaker.DecorateGet(breaker, func() interface{} { return nil })
	circuitbreaker.DecorateAccept(breaker, func(interface{}) {})
	circuitbreaker.DecorateApply(breaker, func(_ interface{}) interface{} { return nil })
}
