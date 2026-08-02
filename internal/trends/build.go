package trends

import "math"

/*
BuildTrends returns a Trends with global trends of values, ie.
the most relevant values from values describing values trends
Each value is the latest value following the same trend (increase or decrease)
If the change is not sensitive (based on threshold), it is ignored

Threshold values example :
- 30 is better for large overview
- 1-6 is better for smoothing elevation gain/loss calculation

The first and last element are always added to the results.
The last element is either the last trend, or a new trend, that both has to be added
*/
func BuildTrends(values []float64, threshold float64) Trends {
	var isHigher, prevIsHigher bool
	var prevVal float64
	var out Trends

	for i, v := range values {
		diff := math.Round(v - prevVal)
		isHigher = (diff > 0)

		// Init first value
		if i == 0 {
			out.add(TrendItem{Index: i, Value: v})
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
					out.update(len(out)-1, TrendItem{Index: i, Value: v})
				}

				// If change is sensitive, so append prevVal to out
			} else {
				out.add(TrendItem{Index: i - 1, Value: prevVal})
				prevIsHigher = isHigher // register the trend change
			}

		}

		// Append last value to out
		// (either the trend is the same, then the last trend is not added yet
		// or the trend is changing on the last item, then it also has to be added
		// so it always need to be added)
		if i == len(values)-1 {
			out.add(TrendItem{Index: i, Value: v})
		}

		prevVal = v
	}

	return out
}
