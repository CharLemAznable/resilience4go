package utils_test

import (
	"github.com/CharLemAznable/resilience4go/utils"
	"testing"
)

func TestMin(t *testing.T) {
	tests := []struct {
		x, y, want int64
	}{
		{3, 5, 3},
		{5, 3, 3},
		{0, 0, 0},
		{-5, -3, -5},
		{-3, -5, -5},
		{9223372036854775807, 9223372036854775807, 9223372036854775807},
		{-9223372036854775808, -9223372036854775808, -9223372036854775808},
		{9223372036854775807, -9223372036854775808, -9223372036854775808},
		{-9223372036854775808, 9223372036854775807, -9223372036854775808},
	}

	for _, test := range tests {
		got := utils.Min(test.x, test.y)
		if got != test.want {
			t.Errorf("Min(%d, %d) = %d; want %d", test.x, test.y, got, test.want)
		}
	}
}
