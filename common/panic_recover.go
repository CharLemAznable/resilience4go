package common

import "fmt"

type Panicked chan any

func (pr Panicked) Recover() {
	if err := recover(); err != nil {
		pr <- err
	}
}

func (pr Panicked) Caught() <-chan any {
	return pr
}

func PanicError(v any) error {
	return &panicError{error: v}
}

type panicError struct {
	error any
}

func (e *panicError) Error() string {
	return fmt.Sprintf("panicked with %v", e.error)
}
