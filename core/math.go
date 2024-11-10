package core

import (
	"fmt"
	"math"
	"slices"
	"strconv"
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

func Dist(p1 Pos, p2 Pos) float64 {
	// output: km
	var R float64 = 6371

	theta2 := DegToRad(p2.Lat)
	theta1 := DegToRad(p1.Lat)
	phi2 := DegToRad(p2.Lon)
	phi1 := DegToRad(p1.Lon)

	h := Haversin(theta2 - theta1)
	h += math.Cos(theta1) * math.Cos(theta2) * Haversin(phi2-phi1)
	out := R * Ahaversin(h)
	return out

}

func DiffElevation(p1 Pos, p2 Pos) float64 {
	return p2.Ele - p1.Ele
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

type RollCalc func([]float64) float64

func Mean(v []float64) float64 {
	var out float64 = 0
	for _, e := range v {
		out += e
	}
	return out / float64(len(v))
}

func Rolling(v []float64, winSize int, calc RollCalc) []float64 {
	var out []float64

	n := len(v)
	// if n < winSize ?

	for i, _ := range v {
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

type IndexValue struct {
	Index int
	Value float64
}

type IndexValues []IndexValue

func VariationSummary(s []float64) IndexValues {
	var lasthigher bool
	var last float64
	var out IndexValues

	for i, v := range s {
		// Append first value
		if i == 0 {
			last = v
			out = append(out, IndexValue{Index: i, Value: v})
			continue
		}

		diff := math.Round(v - last)

		// Init "last" (previous) values
		if i == 1 {
			last = v
			lasthigher = diff > 0
			continue
		}
		// Append last value
		if i == len(s)-1 {
			out = append(out, IndexValue{Index: i, Value: v})
			break
		}

		higher := (diff > 0)
		if higher != lasthigher {

			if math.Abs(out[len(out)-1].Value-last) < 30 && len(out) > 1 {
				if (higher && v > out[len(out)-1].Value) ||
					(!higher && v < out[len(out)-1].Value) {
					out[len(out)-1] = IndexValue{Index: i, Value: v}
				}
			} else {
				out = append(out, IndexValue{Index: i - 1, Value: last})
				lasthigher = higher
			}

		}
		last = v

	}

	return out
}

func (s IndexValues) MinMax() (float64, float64) {
	var l []float64
	for _, v := range s {
		l = append(l, v.Value)
	}
	return slices.Min(l), slices.Max(l)
}
func (s IndexValues) Min() float64 {
	var l []float64
	for _, v := range s {
		l = append(l, v.Value)
	}
	return slices.Min(l)
}

func Printer(s IndexValues) {
	m, M := s.MinMax()

	var lines [5]string
	var space string = "    "

	var prev float64
	for i, v := range s {
		linenum := len(lines) - 1 - int((v.Value-m)/(M-m)*float64(len(lines)-1))
		for j, _ := range lines {
			lines[j] = lines[j] + space
			if j == linenum {
				lines[j] = lines[j] + strconv.Itoa(int(v.Value))
			} else {
				lines[j] = lines[j] + space
			}
		}
		fmt.Printf("[%v]", linenum)
		if i > 0 {
			if v.Value > prev {
				fmt.Printf("/")
			} else {
				fmt.Printf("\\")
			}
		}
		prev = v.Value
		fmt.Printf(" %.0f m ", v.Value)
	}

	fmt.Println()
	fmt.Println("============")
	for _, line := range lines {
		fmt.Println(line)
	}
	fmt.Println("============")
}
