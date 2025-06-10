package core

import (
	"encoding/xml"
	"fmt"
	"os"
	"slices"
)

// TODO: add return *Gpx without testing !
func (gpx *Gpx) ParseFile(gpxFilename string) *Gpx {
	data, _ := os.ReadFile(gpxFilename)
	if err := xml.Unmarshal(data, &gpx); err != nil {
		if err.Error() != "EOF" {
			fmt.Println(err)
		}
	}
	return gpx
}

func (gpx *Gpx) SetVitesse(v float64) {
	gpx.Extensions.Vitesse = v
}

func (gpx Gpx) GetInfo(ascii_format bool) GpxSummary {
	var trkSummary GpxSummary
	for i, trk := range gpx.Trk {
		summary := trk.GetInfo(i, gpx.Extensions.Vitesse, true)
		trkSummary = append(trkSummary, summary)
	}
	return trkSummary
}

func (gpx *Gpx) GetClosestTrkpts(p Pos) []*Trkpt {
	var trkpts []*Trkpt
	// var ind struct{ i, j, k int }

	seg := gpx.Trk[0].Trkseg[0]
	p0 := seg.Trkpt[0]
	minDist := Dist(
		// TODO: add elevation ?
		Pos{Lat: p.Lat, Lon: p.Lon},
		Pos{Lat: p0.Lat, Lon: p0.Lon},
	)

	for i, _ := range gpx.Trk {
		for j, _ := range gpx.Trk[i].Trkseg {
			for k, trkpt := range gpx.Trk[i].Trkseg[j].Trkpt {
				d := Dist(
					Pos{Lat: p.Lat, Lon: p.Lon},
					Pos{Lat: trkpt.Lat, Lon: trkpt.Lon},
				)

				if d == minDist {
					trkpts = append(trkpts, &gpx.Trk[i].Trkseg[j].Trkpt[k])
				} else if d < minDist {
					// Using index to prevent copy value to keep correct address
					trkpts = []*Trkpt{&gpx.Trk[i].Trkseg[j].Trkpt[k]}
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

	slices.Reverse(gpx.Trk)
	for _, trk := range gpx.Trk {
		trk.Reverse()
	}

	return gpx
}

func (gpx Gpx) SplitAtName(name string) Gpx {
	found := false

out:
	for i, trk := range gpx.Trk {
		for j, trkseg := range trk.Trkseg {
			for k, trkpt := range trkseg.Trkpt {
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

func (gpx *Gpx) AddWpt(wpt Wpt) Gpx {
	gpx.Wpt = append(gpx.Wpt, wpt)
	return *gpx
}

// TODO: Split should put trkptid to the previous trk (update: NO)
// TODO: Split should be inplace (method of *Gpx)
// Split Trk[trkId] containing trkptId into two trk 0:trptkId and trkptId:end
func (gpx Gpx) Split(trkId, trksegId, trkptId int) Gpx {
	// filterBeforeTrkpt returns Trk keeping everything BEFORE TrkptId (excluded)
	// return gpx.Trk[trkId].Trkseg[:trkseg+1] where Trkseg[trkSeg] is filter on Trkpt[:trkptEnd]
	filterBeforeTrkpt := func(gpx Gpx, trkId, trksegId, trkptId int, name string) Trk {
		trk := gpx.Trk[trkId]
		trksegs_bef := trk.Trkseg[:trksegId]
		trkseg_last := Trkseg{trk.Trkseg[trksegId].Trkpt[:trkptId]}

		// Create result Trk
		out := Trk{Name: name}
		if len(trksegs_bef) > 0 { // TODO: check
			// if trksegId > 0 { // if non-empty
			out.Trkseg = trksegs_bef
		}
		if len(trkseg_last.Trkpt) > 0 { // TODO: check
			// if trkptId > 0 { // if non-empty
			// out.Trkseg = slices.Concat(out.Trkseg, trksegs_aft)
			out.Trkseg = append(out.Trkseg, trkseg_last)
		}

		return out
	}

	// filterAfterTrkpt returns Trk keeping everything AFTER TrkptId (excluded)
	// return gpx.Trk[trkId].Trkseg[trkseg:] where Trkseg[0] is filter on Trkpt[trkptId:]
	filterAfterTrkpt := func(gpx Gpx, trkId, trksegId, trkptId int, name string) Trk {
		// OTHER SYNTAX
		// ==============
		// // Output everything after trksegId
		// out := Trk{
		// 	Trkseg: gpx.Trk[trkId].Trkseg[trksegId:],
		// 	Name:   name}
		// // Update first trkseg to filter everything after trkptId
		// out.Trkseg[0] = Trkseg{out.Trkseg[0].Trkpt[trkptId:]}
		// ==============

		trk := gpx.Trk[trkId]
		trksegs_after := trk.Trkseg[trksegId:]
		trkseg_first := Trkseg{trk.Trkseg[trksegId].Trkpt[trkptId:]}

		// Output everything after trksegId
		out := Trk{
			Trkseg: trksegs_after,
			Name:   name}
		// Update first trkseg to filter everything after trkptId
		out.Trkseg[0] = trkseg_first

		return out
	}

	bef := filterBeforeTrkpt(gpx, trkId, trksegId, trkptId, gpx.Trk[trkId].Name)
	aft := filterAfterTrkpt(gpx, trkId, trksegId, trkptId, *gpx.Trk[trkId].Trkseg[trksegId].Trkpt[trkptId].Name)

	// Update gpx.Trk value
	newTrk := slices.Clone(gpx.Trk) // prevent unwanted update on gpx argument
	newTrk = slices.Delete(newTrk, trkId, trkId+1)
	newTrk = slices.Insert(newTrk, trkId, aft) // NOTE: poor readability inserting aft, then bef...
	if len(bef.Trkseg) > 0 {
		newTrk = slices.Insert(newTrk, trkId, bef)
	}

	gpx.Trk = newTrk
	return gpx // newGpx
}

// Merge Trk[trkId2] into Trk[trkId1]
func (gpx *Gpx) Merge(trkId1, trkId2 int) Gpx {
	gpx.Trk[trkId1].Trkseg = slices.Concat(gpx.Trk[trkId1].Trkseg, gpx.Trk[trkId2].Trkseg)
	gpx.Trk = slices.Delete(gpx.Trk, trkId2, trkId2+1)
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
	encoder.Indent("", "\t")

	// Write gpx
	gpx.Filepath = ""
	if err = encoder.Encode(gpx); err != nil {
		fmt.Printf("error: %v\n", err)
		return
	}
}

func (gpx *Gpx) AddColor() *Gpx {
	colors := []string{"8e44ad", "ff5733"}
	for i, _ := range gpx.Trk {
		gpx.Trk[i].Extensions.Line.Xmlns = "http://www.topografix.com/GPX/gpx_style/0/2"
		gpx.Trk[i].Extensions.Line.Color = colors[i%len(colors)]
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
	gpx.ParseFile(gpx.Filepath)

	var out TrknameList
	for i, trk := range gpx.Trk {
		out = append(out, Trkname{Id: i, Name: trk.Name})

		if all {
			for j, trkseg := range trk.Trkseg {
				for k, trkpt := range trkseg.Trkpt {
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
