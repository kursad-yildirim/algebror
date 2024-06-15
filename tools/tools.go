package tools

import (
	"math/rand"
	"strconv"
)

func GetRandomInteger(r, d int) int {
	n := rand.Intn(r * d)
	if n == 0 {
		return n + 1
	} else {
		return n
	}
}

func GetRandomDecimal(r, d int) float64 {
	precision := 0
	if d > 4 {
		precision = 100
	} else {
		precision = 10
	}
	n := float64(rand.Intn(r*precision*d)) / float64(precision)
	if n == 0 {
		return n + 1/float64(precision)
	} else {
		return n
	}
}

func FormatNumber(num float64, p int) string {
	return strconv.FormatFloat(num, 'f', p, 64)
}
