package core

import (
	"fmt"
	"testing"
)

func TestMinMax(t *testing.T) {
	s := NewSeries(
		[]int{0, 0, 0, 0, 0},
		[]float64{1, 2, 3, 4, 5})

	have1, have2 := s.MinMax()
	want1, want2 := 1.0, 5.0

	if want1 != have1 {
		t.Fatalf("Have min: %v, Want: %v", have1, want1)
	}
	if want2 != have2 {
		t.Fatalf("Have max: %v, Want: %v", have2, want2)
	}
}

func TestTrendSummary2(t *testing.T) {
	var tests = []struct {
		have, want Series
	}{
		{
			TrendSummary([]float64{100, 200, 300, 200, 100}, 30),
			NewSeries([]int{0, 2, 4},
				[]float64{100, 300, 100}),
		}, {
			TrendSummary([]float64{100, 200, 100, 200, 100}, 30),
			NewSeries([]int{0, 1, 2, 3, 4},
				[]float64{100, 200, 100, 200, 100}),
		}, {
			// must ignore change
			TrendSummary([]float64{100, 200, 180, 200, 100}, 30),
			NewSeries([]int{0, 1, 4},
				[]float64{100, 200, 100}),
		}, {
			// must update value
			TrendSummary([]float64{100, 200, 180, 220, 100}, 30),
			NewSeries([]int{0, 3, 4},
				[]float64{100, 220, 100}),
		},
	}

	for _, test := range tests {
		testname := fmt.Sprintf("Testing...\nHave: %v\nWant: %v\n", test.have, test.want)
		n := len(test.have)
		m := len(test.want)

		t.Run(testname, func(t *testing.T) {
			if n != m {
				t.Fatalf("have length (%v) differs from want length (%v)", n, m)
			}
			for i := range test.have {
				idx1 := test.have[i].Index
				idx2 := test.want[i].Index
				val1 := test.have[i].Value
				val2 := test.want[i].Value

				if idx1 != idx2 {
					t.Fatalf("Items #%v differs: have index (%v) differs from want index (%v)", i, idx1, idx2)
				}
				if val1 != val2 {
					t.Fatalf("Items #%v differs: have value (%v) differs from want value (%v)", i, val1, val2)
				}
			}
		})
	}
}

// func Benchmark1(b *testing.B) {
// 	r := NewSeries(
// 		[]int{0, 0, 0, 0, 0},
// 		[]float64{1, 2, 3, 4, 5})
// 	for i := 0; i < b.N; i++ {
// 		r.MinMax()
// 	}
// }

// func Benchmark2(b *testing.B) {
// 	r := NewSeries(
// 		[]int{0, 0, 0, 0, 0},
// 		[]float64{1, 2, 3, 4, 5})
// 	for i := 0; i < b.N; i++ {
// 		r.MinMax2()
// 	}
// }
