package ratelimiter

type Metrics interface {
	NumberOfWaitingThreads() int64
	AvailablePermissions() int64
}

func newMetrics(numberOfWaitingThreadsFn func() int64, availablePermissionsFn func() int64) *metrics {
	return &metrics{
		numberOfWaitingThreadsFn: numberOfWaitingThreadsFn,
		availablePermissionsFn:   availablePermissionsFn,
	}
}

type metrics struct {
	numberOfWaitingThreadsFn func() int64
	availablePermissionsFn   func() int64
}

func (m *metrics) NumberOfWaitingThreads() int64 {
	return m.numberOfWaitingThreadsFn()
}

func (m *metrics) AvailablePermissions() int64 {
	return m.availablePermissionsFn()
}
