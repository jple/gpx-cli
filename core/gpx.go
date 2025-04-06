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
	for _, trk := range gpx.Trk {
		summary := trk.GetInfo(gpx.Extensions.Vitesse, false)
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

func (gpx Gpx) Save(filepath string) {
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

type TrkName struct {
	TrkName    string
	TrkptNames []string
}
type TrkNames []TrkName

func (gpx Gpx) Ls(all bool) TrkNames {
	gpx.ParseFile(gpx.Filepath)

	var out TrkNames
	for i, trk := range gpx.Trk {
		out = append(out, TrkName{TrkName: trk.Name})

		if all {
			trkpts := slices.Concat(trk.Trkseg)[0].Trkpt
			for _, trkpt := range trkpts {
				if trkpt.Name != nil {
					out[i].TrkptNames = append(out[i].TrkptNames, *trkpt.Name)
				}
			}
		}
	}

	return out
}

func (trkNames TrkNames) Print(all bool, ascii_format ...bool) {

	for i, trkName := range trkNames {
		if len(ascii_format) > 0 && !ascii_format[0] {
			fmt.Printf("[%v] %v\n", i, trkName.TrkName)
		} else {
			fmt.Printf("[%v] \u001b[1;32m%v\u001b[22;0m\n", i, trkName.TrkName)
		}
		if all {
			for _, trkptName := range trkName.TrkptNames {
				fmt.Println(trkptName)
			}
			fmt.Println()
		}
	}
}
