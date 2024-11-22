package core

import (
	"fmt"
	"math"
	"slices"
	"strconv"
	"strings"
)

type (
	IndexValue struct {
		Index int
		Value float64
	}

	Series []IndexValue
)

func (s Series) GetValues() []float64 {
	var values []float64
	for _, v := range s {
		values = append(values, v.Value)
	}
	return values
}
func (s Series) MinMax() (float64, float64) {
	var l []float64
	for _, v := range s {
		l = append(l, v.Value)
	}
	return slices.Min(l), slices.Max(l)
}
func (s Series) Min() float64 {
	var l []float64
	for _, v := range s {
		l = append(l, v.Value)
	}
	return slices.Min(l)
}

// TrendSummary returns a list of Series ({Index, Values])
// The first element is the first value input s[0]
// The next one is the latest increasing (or decreasing) values
// The next one is the latest decréasing (or increasing) values
// And so one until the last element
func TrendSummary(s []float64) Series {
	var isHigher bool
	var prevIsHigher bool
	var prevVal float64
	// Out contains most relevant values from s describing s trends
	// Each value is the latest value following the same trend (increase or decrease)
	// If the change is not sensitive, it is ignored
	var out Series

	for i, v := range s {
		// Append first value
		if i == 0 {
			prevVal = v
			out = append(out, IndexValue{Index: i, Value: v})
			continue
		}

		diff := math.Round(v - prevVal)

		// Init "prevVal" (previous) values
		if i == 1 {
			prevVal = v
			prevIsHigher = diff > 0
			continue
		}
		// Append last value to out
		if i == len(s)-1 {
			out = append(out, IndexValue{Index: i, Value: v})
			break
		}

		isHigher = (diff > 0)

		// Curve is changing dynamics
		if isHigher != prevIsHigher {

			// If latest "out" value and prevVal is too small
			// (ignore first value)
			// The trend is ignored, and the value is not appended in "out"
			// However, updating the latest "out" value may be needed, if the value
			// is higher (in increase trend) than the last "out" value. Same if v
			// is lower (in decrease trend)
			if math.Abs(out[len(out)-1].Value-prevVal) < 30 && len(out) > 1 {

				// Update latest registered values if v is following the same trend
				if (isHigher && v > out[len(out)-1].Value) ||
					(!isHigher && v < out[len(out)-1].Value) {
					out[len(out)-1] = IndexValue{Index: i, Value: v}
				}

				// Otherwise, append prevVal to out
			} else {
				out = append(out, IndexValue{Index: i - 1, Value: prevVal})
				prevIsHigher = isHigher // register the trend change
			}

		}
		prevVal = v
	}

	return out
}

func floor(x float64) int {
	return int(math.Floor(math.Log10(x)) + 1)
}

func replaceAtIndex(str string, replacement rune, index int) string {
	return str[:index] + string(replacement) + str[index+1:]
}

// Print the trend summary based on TrendSummary
// Prints on 5 lines, index 0 for highest values, 4 for lowest ones
func (s Series) PrintTrends() {
	m, M := s.MinMax()

	var space string = strings.Repeat(" ", floor(M))
	var lines [5]string

	var prevLine int
	var prevVal float64
	for i, v := range s {
		// lineToPrint is the index on which the current value is to print
		// It is computed from v compared to m and M
		// The closest to M is to print on 0 index
		// The closest to m is to print on higher index (len(lines) - 1
		lineToPrint := int(len(lines) - 1 - int((v.Value-m)/(M-m)*float64(len(lines)-1)))

		// Prints the value on the corresponding line
		// Including trend information
		if i > 0 {
			if v.Value >= prevVal {
				lines[lineToPrint] = lines[lineToPrint] + replaceAtIndex(space, '/', len(space)-1) + strconv.Itoa(int(v.Value))
			} else {
				lines[lineToPrint] = lines[lineToPrint] + replaceAtIndex(space, '\\', len(space)-1) + strconv.Itoa(int(v.Value))
			}
		} else {
			lines[lineToPrint] = lines[lineToPrint] + space + strconv.Itoa(int(v.Value))
		}

		// Loops over lines to prints spaces
		// Including trends information
		for j := 0; j < len(lines); j++ {
			if j == lineToPrint {
				continue
			}

			// If the line is between the lastest print and the current
			// Prints chars to links the two prints on these lines
			l := int(math.Min(float64(prevLine), float64(lineToPrint)))
			L := int(math.Max(float64(prevLine), float64(lineToPrint)))
			linesToPrintTrends := j >= l && j <= L
			if i > 0 && linesToPrintTrends {
				if v.Value > prevVal {
					lines[j] = lines[j] + replaceAtIndex(space, '/', int(math.Max(float64(prevLine), float64(lineToPrint)))-j) + space
				} else {
					lines[j] = lines[j] + replaceAtIndex(space, '\\', j) + space
				}
				// Otherwise, prints spaces
			} else {
				lines[j] = lines[j] + space + space
			}
		}

		prevLine = lineToPrint
		prevVal = v.Value

	}

	for _, line := range lines {
		fmt.Println(line)
	}
}
