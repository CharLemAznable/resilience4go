package common

import (
	"github.com/CharLemAznable/gofn/common"
)

func Zero[T any]() T {
	return CastOrZero[T](nil)
}

func CastOrZero[T any](val any) T {
	v, _ := common.CastOrZero[T](val)
	return v
}
