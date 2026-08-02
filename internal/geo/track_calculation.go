package geo

import "math"

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
func TotalAscent(s []float64) float64 {
	return SumFunc(s, func(curr, prev float64) float64 {
		return math.Max(curr-prev, 0)
	})
}
func TotalDescent(s []float64) float64 {
	return SumFunc(s, func(curr, prev float64) float64 {
		return math.Min(curr-prev, 0)
	})
}
