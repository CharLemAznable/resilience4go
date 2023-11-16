package utils

import (
	"reflect"
)

func AppendElementUnique[T any](slice []T, value T) []T {
	var removed []T = RemoveElementByValue(slice, value)
	return append(removed, value)
}

func RemoveElementByValue[T any](slice []T, value T) []T {
	result := make([]T, 0)
	start, cursor := 0, 0
	for ; cursor < len(slice); cursor++ {
		if EqualsPointer(slice[cursor], value) {
			if start != cursor {
				result = append(result, slice[start:cursor]...)
			}
			start = cursor + 1
		}
	}
	return append(result, slice[start:cursor]...)
}

func EqualsPointer[T any](x, y T) bool {
	return reflect.ValueOf(x).Pointer() == reflect.ValueOf(y).Pointer()
}

func ConsumeEvent[T any](consumers []func(T), event T) {
	for _, consumer := range consumers {
		go consumer(event)
	}
}
