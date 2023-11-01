package utils

func Min(x, y int64) int64 {
	if x > y {
		return y
	}
	return x
}

func DivCeil(x, y int64) int64 {
	return (x + y - 1) / y
}
