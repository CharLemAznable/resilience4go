package common

import (
	"github.com/CharLemAznable/gofn/common"
)

func Zero[T any]() T {
	zero, _ := common.CastOrZero[T](nil)
	return zero
}
