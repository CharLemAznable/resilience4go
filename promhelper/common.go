package promhelper

import "github.com/prometheus/client_golang/prometheus"

const (
	labelKeyName  = "name"
	labelKeyKind  = "kind"
	labelKeyState = "state"
)

type RegisterFn func(prometheus.Registerer) error

type UnregisterFn func(prometheus.Registerer) bool

func buildRegisterFn(collectors ...prometheus.Collector) RegisterFn {
	return func(registerer prometheus.Registerer) error {
		err := prometheus.MultiError{}
		for _, collector := range collectors {
			err.Append(registerer.Register(collector))
		}
		return err.MaybeUnwrap()
	}
}

func buildUnregisterFn(collectors ...prometheus.Collector) UnregisterFn {
	return func(registerer prometheus.Registerer) bool {
		ret := true
		for _, collector := range collectors {
			ret = ret && registerer.Unregister(collector)
		}
		return ret
	}
}
