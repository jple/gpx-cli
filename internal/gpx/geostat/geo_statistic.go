package geostat

const DefaultIndent int = 4

// GeoStatistic summaries either
// - the whole trk
// - a section between a trkpt name and the next one (no matter trkseg)
type GeoStatistic struct {
	FlatSpeed float64 `json:",omitempty"`

	// Trk.Name or Trkpt.Name (dupplicates with From)
	Name string `json:"name,omitempty"`

	// Trk.Name or Trkpt.Name
	From string `json:",omitempty"`
	To   string `json:",omitempty"`

	// Cumulative values between "From" and "To"
	N              int     `json:",omitempty"` // N is the trkpts count
	Distance       float64 `json:",omitempty"`
	TotalAscent    float64 `json:",omitempty"`
	TotalDescent   float64 `json:",omitempty"`
	DistanceEffort float64 `json:",omitempty"`
	DurationHour   int8    `json:",omitempty"`
	DurationMin    int8    `json:",omitempty"`
}

func (s GeoStatistic) Len() int {
	return s.N
}
