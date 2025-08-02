package core

import (
	"encoding/xml"
	"fmt"
	"os"
	"slices"
	"strings"
)

type Gpx struct {
	XMLName string `xml:"gpx"`
	// NOTE: known issues: xmlns:_xmlns not found
	Attrs []xml.Attr `xml:",any,attr"`

	Trks []Trk `xml:"trk,omitempty"`
	Wpts []Wpt `xml:"wpt,omitempty"`

	Metadata *struct {
		Inner string `xml:",innerxml"`
	} `xml:"metadata,omitempty"`

	Extensions *struct {
		Inner string `xml:",innerxml"`
	} `xml:"extensions,omitempty"`
}

func (gpx *Gpx) ParseFile(gpxFilename string) *Gpx {
	// Parse gpx
	data, _ := os.ReadFile(gpxFilename)
	if err := xml.Unmarshal(data, &gpx); err != nil {
		if err.Error() != "EOF" {
			fmt.Println(err)
		}
	}

	// Cleaning struct
	// ... tags ",innerxml": remove empty struct
	if gpx.Metadata != nil && strings.TrimSpace(gpx.Metadata.Inner) == "" {
		gpx.Metadata = nil
	}
	if gpx.Extensions != nil && strings.TrimSpace(gpx.Extensions.Inner) == "" {
		gpx.Extensions = nil
	}
	// ... tags ",any": remove xmlns attribute
	for i, trk := range gpx.Trks {
		if trk.Extensions != nil {
			for j := range trk.Extensions.Else {
				gpx.Trks[i].Extensions.Else[j].XMLName.Space = ""
			}
		}
	}
	return gpx
}

// func (gpx *Gpx) SetVitesse(v float64) {
// 	gpx.Extensions.Vitesse = v
// }

func (gpx Gpx) GetInfo(vitessePlat float64) GpxSummary {
	var gpxSummary GpxSummary
	for i, trk := range gpx.Trks {
		trkSummary := trk.GetInfo(i, vitessePlat)
		gpxSummary = append(gpxSummary, trkSummary)
	}
	return gpxSummary
}

// Returns slice of pointers to Trkpt
// This function is used to add name to trkpt in-place
// without need to specify Trk,Trkseg,Trkpt id
func (gpx *Gpx) GetClosestTrkpts(p Pt) []*Trkpt {
	var trkpts []*Trkpt
	// var ind struct{ i, j, k int }

	// TODO: add check on len(...)
	seg := gpx.Trks[0].Trksegs[0]
	p0 := seg.Trkpts[0]
	minDist := Dist(
		Pt{Lat: p.Lat, Lon: p.Lon},
		Pt{Lat: p0.Lat, Lon: p0.Lon},
	)

	for i, _ := range gpx.Trks {
		for j, _ := range gpx.Trks[i].Trksegs {
			for k, trkpt := range gpx.Trks[i].Trksegs[j].Trkpts {
				d := Dist(
					Pt{Lat: p.Lat, Lon: p.Lon},
					Pt{Lat: trkpt.Lat, Lon: trkpt.Lon},
				)

				if d == minDist {
					trkpts = append(trkpts, &gpx.Trks[i].Trksegs[j].Trkpts[k])
				} else if d < minDist {
					// Using index to prevent copy value to keep correct address
					trkpts = []*Trkpt{&gpx.Trks[i].Trksegs[j].Trkpts[k]}
					// ind = struct{ i, j, k int }{i, j, k}
					minDist = d
				}
			}
		}
	}

	return trkpts
}

func (p_gpx *Gpx) Reverse() Gpx {
	gpx := *p_gpx

	slices.Reverse(gpx.Trks)
	for _, trk := range gpx.Trks {
		trk.Reverse()
	}

	return gpx
}

func (gpx *Gpx) AddWpt(wpt Wpt) Gpx {
	gpx.Wpts = append(gpx.Wpts, wpt)
	return *gpx
}

// TODO: Split should be inplace (method of *Gpx)
// Split Trk[trkId] containing trkptId into two trk 0:trptkId and trkptId:end
func (gpx Gpx) Split(trkId, trksegId, trkptId int) Gpx {
	// filterBeforeTrkpt returns Trk keeping everything BEFORE TrkptId (excluded)
	// return gpx.Trks[trkId].Trksegs[:trkseg+1] where Trkseg[trkSeg] is filter on Trkpt[:trkptEnd]
	filterBeforeTrkpt := func(gpx Gpx, trkId, trksegId, trkptId int, name string) Trk {
		trk := gpx.Trks[trkId]
		trksegs_bef := trk.Trksegs[:trksegId]
		trkseg_last := Trkseg{trk.Trksegs[trksegId].Trkpts[:trkptId]}

		// Create result Trk
		out := Trk{Name: name}
		if len(trksegs_bef) > 0 { // TODO: check
			// if trksegId > 0 { // if non-empty
			out.Trksegs = trksegs_bef
		}
		if len(trkseg_last.Trkpts) > 0 { // TODO: check
			// if trkptId > 0 { // if non-empty
			// out.Trksegs = slices.Concat(out.Trksegs, trksegs_aft)
			out.Trksegs = append(out.Trksegs, trkseg_last)
		}

		return out
	}

	// filterAfterTrkpt returns Trk keeping everything AFTER TrkptId (excluded)
	// return gpx.Trks[trkId].Trksegs[trkseg:] where Trkseg[0] is filter on Trkpt[trkptId:]
	filterAfterTrkpt := func(gpx Gpx, trkId, trksegId, trkptId int, name string) Trk {
		// OTHER SYNTAX
		// ==============
		// // Output everything after trksegId
		// out := Trk{
		// 	Trkseg: gpx.Trks[trkId].Trksegs[trksegId:],
		// 	Name:   name}
		// // Update first trkseg to filter everything after trkptId
		// out.Trksegs[0] = Trkseg{out.Trksegs[0].Trkpts[trkptId:]}
		// ==============

		trk := gpx.Trks[trkId]
		trksegs_after := trk.Trksegs[trksegId:]
		trkseg_first := Trkseg{trk.Trksegs[trksegId].Trkpts[trkptId:]}

		// Output everything after trksegId
		out := Trk{
			Trksegs: trksegs_after,
			Name:    name}
		// Update first trkseg to filter everything after trkptId
		out.Trksegs[0] = trkseg_first

		return out
	}

	bef := filterBeforeTrkpt(gpx, trkId, trksegId, trkptId, gpx.Trks[trkId].Name)
	aft := filterAfterTrkpt(gpx, trkId, trksegId, trkptId, *gpx.Trks[trkId].Trksegs[trksegId].Trkpts[trkptId].Name)

	// Update gpx.Trks value
	newTrk := slices.Clone(gpx.Trks) // prevent unwanted update on gpx argument
	newTrk = slices.Delete(newTrk, trkId, trkId+1)
	newTrk = slices.Insert(newTrk, trkId, aft) // NOTE: poor readability inserting aft, then bef...
	if len(bef.Trksegs) > 0 {
		newTrk = slices.Insert(newTrk, trkId, bef)
	}

	gpx.Trks = newTrk
	return gpx // newGpx
}

func (gpx Gpx) SplitAtName(name string) Gpx {
	found := false

out:
	for i, trk := range gpx.Trks {
		for j, trkseg := range trk.Trksegs {
			for k, trkpt := range trkseg.Trkpts {
				if trkpt.Name != nil && *trkpt.Name == name {
					found = true

					// TODO: this print is a pb in tui module
					// fmt.Printf("Split at trk %v, trkseg %v, trkpt %v\n", i, j, k)

					gpx = gpx.Split(i, j, k)

					break out
				}
			}
		}
	}

	if !found {
		fmt.Printf("Name ('%v') not found in gpx", name)
	}
	return gpx
}

// Merge Trk[trkId2] into Trk[trkId1]
func (gpx *Gpx) Merge(trkId1, trkId2 int) Gpx {
	gpx.Trks[trkId1].Trksegs = slices.Concat(gpx.Trks[trkId1].Trksegs, gpx.Trks[trkId2].Trksegs)
	gpx.Trks = slices.Delete(gpx.Trks, trkId2, trkId2+1)
	return *gpx
}
func (gpx Gpx) Save(filepath string) {
	if filepath == "" {
		filepath = "out.gpx"
	}
	// TODO: to reset, but mess up with TUI...
	// fmt.Println("Save to", filepath)

	// Create xml file
	xmlFile, err := os.Create(filepath)
	if err != nil {
		fmt.Println("Error creating XML file:", err)
		return
	}

	// Write xml header
	_, err = xmlFile.Write([]byte(xml.Header))
	if err != nil {
		fmt.Println("Error writing to XML file:", err)
		return
	}

	encoder := xml.NewEncoder(xmlFile)
	encoder.Indent("", "  ")

	// Write gpx
	if err = encoder.Encode(gpx); err != nil {
		fmt.Printf("error: %v\n", err)
		return
	}
}

func (gpx *Gpx) AddColor() *Gpx {
	colors := []string{"8e44ad", "ff5733"}
	newLineColor := ExtensionsLine{
		Attrs: []xml.Attr{
			xml.Attr{
				xml.Name{"", "xmlns"},
				"http://www.topografix.com/GPX/gpx_style/0/2",
			}},
	}

	for i, _ := range gpx.Trks {
		newLineColor.Color = colors[i%len(colors)]

		// TODO: improvement: create Trk.AddLineColor
		if gpx.Trks[i].Extensions == nil {
			gpx.Trks[i].Extensions = &ExtensionsTrk{Line: &newLineColor}
		} else {
			gpx.Trks[i].Extensions.Line = &newLineColor
		}
	}
	return gpx
}

// TODO: To be removed after update tui/model.go
// ================= Ls and Print ===========================
type (
	Trkname struct {
		Id            int
		Name          string
		TrksegId      *int
		TrkptId       *int
		TrkptName     *string
		Lat, Lon, Ele *float64
	}
	TrknameList []Trkname
)

func (tn Trkname) IsTrkpt() bool {
	if tn.TrkptName != nil {
		return true
	}
	return false
}
func (tn Trkname) IsTrk() bool {
	return !tn.IsTrkpt()
}

func (gpx Gpx) Ls(all bool) TrknameList {
	var out TrknameList
	for i, trk := range gpx.Trks {
		out = append(out, Trkname{Id: i, Name: trk.Name})

		if all {
			for j, trkseg := range trk.Trksegs {
				for k, trkpt := range trkseg.Trkpts {
					if trkpt.Name != nil {
						out = append(out,
							Trkname{
								Id:        i,
								Name:      trk.Name,
								TrksegId:  &j,
								TrkptId:   &k,
								TrkptName: trkpt.Name,
								Lat:       &trkpt.Lat,
								Lon:       &trkpt.Lon,
								Ele:       &trkpt.Ele,
							})
					}
				}
			}
		}
	}

	return out
}

func (tnList TrknameList) Print(all bool, ascii_format ...bool) {

	for _, trkname := range tnList {
		if trkname.IsTrk() {
			if len(ascii_format) > 0 && !ascii_format[0] {
				fmt.Printf("[%v] %v\n", trkname.Id, trkname.Name)
			} else {
				fmt.Printf("[%v] \u001b[1;32m%v\u001b[22;0m\n", trkname.Id, trkname.Name)
			}
		}
		if all {
			// fmt.Printf("(seg:%v, pt:%v) %v\n", pt.TrksegId, pt.Id, pt.Name)
			if trkname.IsTrkpt() {
				fmt.Printf("\t%v %v %v %v\n",
					*trkname.TrkptName,
					*trkname.Lat, *trkname.Lon, *trkname.Ele)
			}
		}
	}
}
