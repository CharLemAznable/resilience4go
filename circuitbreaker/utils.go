package circuitbreaker

import "sync/atomic"

func getAndUpdatePointer[T any](pointer *atomic.Pointer[T], updateFn func(*T) (*T, error)) (*T, error) {
	prev := pointer.Load()
	var next *T
	var err error
	for haveNext := false; ; {
		if !haveNext {
			next, err = updateFn(prev)
			if err != nil {
				return nil, err
			}
		}
		if pointer.CompareAndSwap(prev, next) {
			return prev, nil
		}
		temp, prev := prev, pointer.Load()
		haveNext = temp == prev
	}
}
