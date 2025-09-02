package core

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

type IndexValue struct {
	Index int
	Value float64
}
type Series []IndexValue

func (s *Series) Add(e IndexValue) Series {
	*s = append(*s, e)
	// fmt.Printf("Add new item #%v: %v\n", len(*s)-1, e)
	return *s
}

func (s *Series) Update(i int, e IndexValue) Series {
	(*s)[i] = e
	// fmt.Printf("Update item #%v: %v\n", i, e)
	return *s
}

func NewSeries(idxs []int, values []float64) Series {
	var s Series

	if len(idxs) != len(values) {
		panic("idxs and values don't have the same length")
	}

	for i := range idxs {
		s.Add(IndexValue{idxs[i], values[i]})
	}
	return s
}

func (s Series) GetValues() []float64 {
	var values []float64
	for _, v := range s {
		values = append(values, v.Value)
	}
	return values
}

/*
TrendSummary returns a Series with global trends of s, ie.
the most relevant values from s describing s trends
Each value is the latest value following the same trend (increase or decrease)
If the change is not sensitive (based on threshold), it is ignored

Threshold values example :
- 30 is better for large overview
- 1-6 is better for smoothing elevation gain/loss calculation

The first and last element are always added to the results.
The last element is either the last trend, or a new trend, that both has to be added
*/
func TrendSummary(s []float64, threshold float64) Series {
	var isHigher, prevIsHigher bool
	var prevVal float64
	var out Series

	for i, v := range s {
		diff := math.Round(v - prevVal)
		isHigher = (diff > 0)

		// fmt.Printf(
		// 	"i (%v), v (%v), prev (%v), diff (%v), higher (%v), prevH (%v)\n",
		// 	i, v, prevVal, diff, isHigher, prevIsHigher)

		// Init first value
		if i == 0 {
			out.Add(IndexValue{Index: i, Value: v})
			prevVal = v
			continue
		}

		// Init first trend
		if i == 1 {
			prevIsHigher = diff > 0
			prevVal = v
			continue
		}

		// Curve is changing dynamics
		if isHigher != prevIsHigher {

			lastTrendValue := out[len(out)-1].Value
			// If change in value is too small, it is generally ignored,
			// except if the value is higher (in increasing trend), or lower (in decreasing trend)
			// than the registered one, so the value is updated
			if math.Abs(lastTrendValue-prevVal) < threshold && len(out) > 1 {

				// Update latest registered values if v is following the same trend
				if (isHigher && v > lastTrendValue) ||
					(!isHigher && v < lastTrendValue) {
					// fmt.Printf("previous value #%v: %v\n", len(out)-1, out[len(out)-1].Value)
					out.Update(len(out)-1, IndexValue{Index: i, Value: v})
				}

				// If change is sensitive, so append prevVal to out
			} else {
				out.Add(IndexValue{Index: i - 1, Value: prevVal})
				prevIsHigher = isHigher // register the trend change
			}

		}

		// Append last value to out
		// (either the trend is the same, then the last trend is not added yet
		// or the trend is changing on the last item, then it also has to be added
		// so it always need to be added)
		if i == len(s)-1 {
			out.Add(IndexValue{Index: i, Value: v})
		}

		prevVal = v
	}

	return out
}

// ======
// TODO: move to term-plot ?

func floor(x float64) int {
	return int(math.Floor(math.Log10(x)) + 1)
}

func replaceAtIndex(str string, replacement rune, index int) string {
	return str[:index] + string(replacement) + str[index+1:]
}

// TODO: Must not be useful... except for term-plot...
func (s Series) IndexMinMax() (int, float64, int, float64) {
	m, M := s[0].Value, s[0].Value
	imin, imax := s[0].Index, s[0].Index
	for i := range s {
		if m > s[i].Value {
			imin = s[i].Index
			m = s[i].Value
		}
		if M < s[i].Value {
			imax = s[i].Index
			M = s[i].Value
		}
	}
	return imin, m, imax, M
}
func (s Series) MinMax() (float64, float64) {
	_, m, _, M := s.IndexMinMax()
	return m, M
}
func (s Series) Min() float64 {
	_, m, _, _ := s.IndexMinMax()
	return m
}

// Print the trend summary based on TrendSummary
// Prints on 5 lines, index 0 for highest values, 4 for lowest ones
func (s Series) PrintTrends() {
	m, M := s.MinMax()
	if m == M {
		fmt.Println("No variations in this series: min = max =", m)
		return
	}

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
