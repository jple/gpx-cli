package core

import (
	"encoding/xml"
	"fmt"
	"os"
	"slices"
)

func (gpx *Gpx) ParseFile(gpxFilename string) {
	data, _ := os.ReadFile(gpxFilename)
	if err := xml.Unmarshal(data, &gpx); err != nil {
		if err.Error() != "EOF" {
			fmt.Println(err)
		}
	}
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
		fmt.Println("Name (", name, ") not found in gpx")
	}
	return gpx
}

func (gpx *Gpx) AddWpt(wpt Wpt) Gpx {
	gpx.Wpt = append(gpx.Wpt, wpt)
	return *gpx
}

// TODO: split should put trkptid to the previous trk
func (gpx Gpx) Split(trkId, trksegId, trkptId int) Gpx {
	filterBeforeTrkpt := func(gpx Gpx, trkId, trksegId, trkptId int, name string) Trk {
		// return gpx.Trk[trkId].Trkseg[:trkseg+1] where Trkseg[trkSeg] is filter on Trkpt[:trkptEnd]

		trk := gpx.Trk[trkId]
		trk.AddName(name)

		trk.Trkseg = slices.Concat(
			trk.Trkseg[:trksegId],
			[]Trkseg{
				Trkseg{trk.Trkseg[trksegId].Trkpt[:trkptId]},
			},
		)
		return trk
	}
	filterAfterTrkpt := func(gpx Gpx, trkId, trksegId, trkptId int, name string) Trk {
		// return gpx.Trk[trkId].Trkseg[trkseg:] where Trkseg[0] is filter on Trkpt[trkptId:]

		trk := gpx.Trk[trkId]
		trk.AddName(name)

		trk.Trkseg = trk.Trkseg[trksegId:]
		trk.Trkseg[0] = Trkseg{trk.Trkseg[0].Trkpt[trkptId:]}
		return trk
	}

	// WIP
	bef := filterBeforeTrkpt(gpx, trkId, trksegId, trkptId, gpx.Trk[trkId].Name)
	aft := filterAfterTrkpt(gpx, trkId, trksegId, trkptId, *gpx.Trk[trkId].Trkseg[trksegId].Trkpt[trkptId].Name)

	gpx.Trk = slices.Delete(gpx.Trk, trkId, trkId+1)
	gpx.Trk = slices.Insert(
		gpx.Trk, trkId,
		bef, aft,
	)

	return gpx
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
	fmt.Println("Save to", filepath)

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
