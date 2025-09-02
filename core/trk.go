package core

import (
	"encoding/xml"
	"slices"
	"strconv"
)

type Trkseg struct {
	Trkpts Trkpts `xml:"trkpt"`
}

type ExtensionsTrk struct {
	Line *ExtensionsLine `xml:"line,omitempty"`
	Else []struct {
		XMLName xml.Name
		Content string     `xml:",innerxml"`
		Attrs   []xml.Attr `xml:",any,attr"`
	} `xml:",any"`
}
type ExtensionsLine struct {
	Attrs []xml.Attr `xml:",any,attr"`

	Color    string `xml:"color,omitempty"`
	Opacity  string `xml:"opacity,omitempty"`
	Weight   string `xml:"Weight,omitempty"`
	Width    int    `xml:"width,omitempty"`
	Linecap  string `xml:"linecap,omitempty"`
	Linejoin string `xml:"linejoin,omitempty"`

	// NOTE: I don't know why, but if Dasharray is not pointer
	// then plot does not work correctly...
	// So Dashoffset has also change to pointer, just in case...
	Dasharray  *int `xml:"dasharray,omitempty"`
	Dashoffset *int `xml:"dashoffset,omitempty"`

	Else []struct {
		XMLName xml.Name
		Content string `xml:",innerxml"`
	} `xml:",any"`
}

type Trk struct {
	Name string `xml:"name,omitempty"`
	// NOTE: innerxml to prevent escaping (more readable, less secure :/)
	// NOTE: I don't know why, but it breaks plot...Probable due to escape issue
	// Name       string         `xml:",innerxml"`
	Extensions *ExtensionsTrk `xml:"extensions,omitempty"`

	Trksegs []Trkseg `xml:"trkseg"`
}

func (trk *Trk) SetName(name string) {
	trk.Name = name
}

func (trk Trk) GetTrkpts() Trkpts {
	trkpts := Trkpts{}
	for _, trkseg := range trk.Trksegs {
		trkpts = slices.Concat(trkpts, trkseg.Trkpts)
	}
	return trkpts
}

func (trk Trk) GetLonLat() ([]string, []string) {
	var lons, lats []string
	for _, trkpt := range trk.GetTrkpts() {
		lons = append(lons, strconv.FormatFloat(trkpt.Lon, 'f', -1, 64))
		lats = append(lats, strconv.FormatFloat(trkpt.Lat, 'f', -1, 64))
	}
	return lons, lats
}

func (p_trk *Trk) Reverse() Trk {
	trk := *p_trk

	slices.Reverse(trk.Trksegs)
	for _, trkseg := range trk.Trksegs {
		slices.Reverse(trkseg.Trkpts)
	}
	return trk
}
