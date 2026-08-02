package gpx

import (
	"encoding/xml"
	"slices"
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

// TrkSections is a list of Trkpts
// This list is a split of the whole Trkpts, in parts of interest
// For example, Trkpts split by named Trkpt
type TrkSections []Trkpts

type Trk struct {
	Name string `xml:"name,omitempty"`
	// NOTE: innerxml to prevent escaping (more readable, less secure :/)
	// NOTE: I don't know why, but it breaks plot...Probable due to escape issue
	// Name       string         `xml:",innerxml"`
	Extensions *ExtensionsTrk `xml:"extensions,omitempty"`

	Trksegs []Trkseg `xml:"trkseg"`
}

func (trk Trk) AllTrkpts() Trkpts {
	trkpts := Trkpts{}
	for _, trkseg := range trk.Trksegs {
		trkpts = slices.Concat(trkpts, trkseg.Trkpts)
	}
	return trkpts
}

func (trk *Trk) Reverse() {
	slices.Reverse(trk.Trksegs)
	for _, trkseg := range trk.Trksegs {
		slices.Reverse(trkseg.Trkpts)
	}
}

// SplitByTrkptName creates TrkSections
// where each Section is a []Trkpt from named item to next named item (included)
// NOTE: the next named item is included in order to correctly calculate geostatistic of Trkpts
// Without it, the calculation would skip the calculation between TrkSections[i].Section[-1] and
// TrkSections[i+1].Section[0]
func (trk Trk) SplitByTrkptName() TrkSections {
	currentSection := Trkpts{}
	TrkSections := TrkSections{}

	for i, trkpt := range trk.AllTrkpts() {
		currentSection = append(currentSection, trkpt)

		// Reach a new named pt (which is not first element)
		// This is the end of the current "section" (currentSection)
		if trkpt.Name != nil && i > 0 {
			// Append "section" to TrkSections
			TrkSections = append(TrkSections, currentSection)
			// Prepare the new section, with the named pt
			currentSection = Trkpts{trkpt}
		}
	}
	// Add last "section" (trkpts) into the list
	TrkSections = append(TrkSections, currentSection)
	return TrkSections
}
