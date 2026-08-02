package geo

import (
	"math"
)

func Haversin(t float64) float64 {
	return math.Pow(math.Sin(t/2), 2)
}

func Ahaversin(t float64) float64 {
	return 2 * math.Asin(math.Sqrt(t))
}

func DegToRad(deg float64) float64 {
	return 2 * math.Pi / 360 * deg
}

func FloatToHourMin(f float64) (int8, int8) {
	hour := math.Floor(f)
	minute := math.Floor((f - hour) * 60)
	return int8(hour), int8(minute)
}

func Sum(v []float64) float64 {
	var out float64 = 0
	for _, e := range v {
		out += e
	}
	return out
}

func Mean(v []float64) float64 {
	return Sum(v) / float64(len(v))
}

func Rolling(v []float64, winSize int, calc func([]float64) float64) []float64 {
	var out []float64

	n := len(v)
	// TODO: if n < winSize ?

	// for i := range v {
	for i := 0; i < len(v)-winSize+1; i++ {
		var s []float64
		for j := 0; j < winSize; j++ {
			s = append(s, v[i+j])
		}
		out = append(out, calc(s))

		if i == n-winSize {
			break
		}
	}

	return out
}

// SumFunc returns sum of f(s[i], s[i-1])
func SumFunc(s []float64, f func(curr, prev float64) float64) float64 {
	var acc, prev float64
	for i, curr := range s {
		if i > 0 {
			acc += f(curr, prev)
		}
		prev = curr
	}
	return acc
}
