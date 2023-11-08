package utils

import (
	"reflect"
)

type Slices[T any] interface {
	RemoveElementByValue(slice []T, value T) []T
	AppendElementUnique(slice []T, value T) []T
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
	result := make([]T, 0)
	start, cursor := 0, 0
	for ; cursor < len(slice); cursor++ {
		if s.equals(slice[cursor], value) {
			if start != cursor {
				result = append(result, slice[start:cursor]...)
			}
			start = cursor + 1
		}
	}
	return append(result, slice[start:cursor]...)
}

func (s *slices[T]) AppendElementUnique(slice []T, value T) []T {
	var removed []T = s.RemoveElementByValue(slice, value)
	return append(removed, value)
}
