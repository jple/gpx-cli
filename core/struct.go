// This file is automatically generated

package core

type Gpx struct {
	XMLName  string `xml:"gpx"`
	Filepath string `xlm:-` // ok for json, but xml ?

	Metadata struct {
		Desc string `xml:"desc"`
		Name string `xml:"name"`
	} `xml:"metadata"`
	Trk []Trk `xml:"trk"`
	Wpt []Wpt `xml:"wpt"`

	Extensions struct {
		Vitesse float64 "Vitesse de marche sur plat (km/h)"
	} `xml:"extensions"`
}

type Trk struct {
	Name string `xml:"name"`
	// Should be optional
	Extensions struct {
		DenivPos       float64
		DenivNeg       float64
		Distance       float64
		DenivPosEffort float64 "Conversion du denivele positif en km effort"
		DenivNegEffort float64 "Conversion du denivele negatif en km effort"
		DistanceEffort float64 "Distance équivalente sur plat en incluant le dénivelé"
		Duration       float64 "Estimation de temps de marche"
		DurationHour   int8
		DurationMin    int8

		Line struct {
			Xmlns      string `xml:"xmlns,attr"`
			Color      string `xml:"color"`
			Dasharray  *int   `xml:"dasharray"`
			Extensions *struct {
				Jonction int `xml:"jonction"`
			} `xml:"extensions"`
			Opacity *float64 `xml:"opacity"`
			Width   int      `xml:"width"`
		} `xml:"line"`
	} `xml:"extensions"`

	Trkseg []Trkseg `xml:"trkseg"`
}

type Trkseg struct {
	Trkpt []Trkpt `xml:"trkpt"`
}

type Trkpt struct {
	Lat float64 `xml:"lat,attr"`
	Lon float64 `xml:"lon,attr"`
	Ele float64 `xml:"ele"`

	Name *string `xml:"name"`
	Type *string `xml:"type"`
	Cmt  *string `xml:"cmt"`

	Extensions *struct {
		TrkExtension struct {
			Visugpx string `xml:"visugpx,attr"`
			Node    int    `xml:"node"`
		} `xml:"TrkExtension"`
	} `xml:"extensions"`
}

type Wpt struct {
	Lat  float64 `xml:"lat,attr"`
	Lon  float64 `xml:"lon,attr"`
	Ele  float64 `xml:"ele"`
	Name *string `xml:"name"`
	Type *string `xml:"type"`
	Cmt  *string `xml:"cmt"`
}
