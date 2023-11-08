package utils

import "reflect"

type Slices[T any] interface {
	RemoveElementByValue(slice []T, value T) []T
}

func NewSlices[T any](equals func(T, T) bool) Slices[T] {
	return &slices[T]{equals: equals}
}

func NewSlicesWithPointer[T any]() Slices[T] {
	return NewSlices[T](func(x, y T) bool {
		return reflect.ValueOf(x).Pointer() == reflect.ValueOf(y).Pointer()
	})
}

type slices[T any] struct {
	equals func(T, T) bool
}

func (s *slices[T]) RemoveElementByValue(slice []T, value T) []T {
	for i := 0; i < len(slice); i++ {
		if s.equals(slice[i], value) {
			result := make([]T, 0)
			return append(append(result, slice[:i]...), slice[i+1:]...)
		}
	}
	return slice
}
