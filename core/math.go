package core

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

func CalcDistanceEffort(distance float64, denivPos float64, denivNeg float64) float64 {
	denivPosEffort := denivPos / 100.0 // 100m D+ = 1km plat
	denivNegEffort := denivNeg / 300.0 // 300m D- = 1km plat
	return distance + denivPosEffort + math.Abs(denivNegEffort)
}

func CalcDuration(distance_plat float64, vitesse_plat float64) (float64, int8, int8) {

	duration := distance_plat / vitesse_plat
	durationHour, durationMin := FloatToHourMin(duration)

	return duration, durationHour, durationMin
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

func CumAscent(s []float64) float64 {
	return SumFunc(s, func(curr, prev float64) float64 {
		return math.Max(curr-prev, 0)
	})
}
func CumDescent(s []float64) float64 {
	return SumFunc(s, func(curr, prev float64) float64 {
		return math.Min(curr-prev, 0)
	})
}
