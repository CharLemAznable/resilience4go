package utils

func Min(x, y int64) int64 {
	if x > y {
		return y
	}
	return x
}

func Abs(x int64) int64 {
	if x < 0 {
		return -x
	}
	return x
}

func DivCeil(x, y int64) int64 {
	if (x < 0 && y > 0) || (x > 0 && y < 0) {
		return -(Abs(x) + Abs(y) - 1) / Abs(y)
	}
	return (Abs(x) + Abs(y) - 1) / Abs(y)
}
