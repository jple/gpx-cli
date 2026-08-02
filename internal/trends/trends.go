package trends

// Trends stores trends calculated from a []float64
// It keeps only turnin point from the slice
// Examples: trends on elevations, cumulative distances
type Trends []TrendItem

// TrendItem is one item of Trends. It is composed of a numerical value (e.g., elevation, distance) at a specific index.
type TrendItem struct {
	Index int     // Index of turning point
	Value float64 // Value at turning point index
}

func NewTrends(index []int, values []float64) Trends {
	var trends Trends

	if len(index) != len(values) {
		panic("index and values don't have the same length")
	}

	for i := range index {
		trends.add(TrendItem{index[i], values[i]})
	}
	return trends
}

func (values *Trends) add(v TrendItem) {
	*values = append(*values, v)
}

func (values *Trends) update(k int, v TrendItem) {
	(*values)[k] = v
}

// func (s Trends) indexMinMax() (int, float64, int, float64) {
func (s Trends) minMax() (float64, float64) {
	m, M := s[0].Value, s[0].Value
	// imin, imax := s[0].Index, s[0].Index
	for i := range s {
		if m > s[i].Value {
			// imin = s[i].Index
			m = s[i].Value
		}
		if M < s[i].Value {
			// imax = s[i].Index
			M = s[i].Value
		}
	}
	// return imin, m, imax, M
	return m, M
}
