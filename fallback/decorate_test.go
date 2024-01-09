package fallback_test

type TargetError struct {
	msg string
}

func (e *TargetError) Error() string {
	return e.msg
}

type NonTargetError struct {
	msg string
}

func (e *NonTargetError) Error() string {
	return e.msg
}
