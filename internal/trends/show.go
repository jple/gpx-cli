package trends

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

func replace(str string, replacement rune, index int) string {
	return str[:index] + string(replacement) + str[index+1:]
}

// Print the trends based on BuildTrends
// Prints on 5 lines, index 0 for highest values, 4 for lowest ones
func (s Trends) ToString() string {
	m, M := s.minMax()
	if m == M {
		return fmt.Sprintln("No variations in this summary: min = max =", m)
	}

	var space string = strings.Repeat(" ", int(math.Floor(math.Log10(M))+1))
	lines := make([]string, 5)

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
				lines[lineToPrint] = lines[lineToPrint] + replace(space, '/', len(space)-1) + strconv.Itoa(int(v.Value))
			} else {
				lines[lineToPrint] = lines[lineToPrint] + replace(space, '\\', len(space)-1) + strconv.Itoa(int(v.Value))
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
			linesToPlotTrends := j >= l && j <= L
			if i > 0 && linesToPlotTrends {
				if v.Value > prevVal {
					lines[j] = lines[j] + replace(space, '/', int(math.Max(float64(prevLine), float64(lineToPrint)))-j) + space
				} else {
					lines[j] = lines[j] + replace(space, '\\', j) + space
				}
				// Otherwise, prints spaces
			} else {
				lines[j] = lines[j] + space + space
			}
		}

		prevLine = lineToPrint
		prevVal = v.Value

	}

	return strings.Join(lines, "\n")
}
