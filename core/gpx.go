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
		fmt.Println(err)
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
		// for _, s := range summary {

		// 	fmt.Printf("[%v] ", i)
		// 	s.Print()
		// }
		// fmt.Println()
	}
	return trkSummary
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
	xmlFile, err := os.Create(filepath)
	if err != nil {
		fmt.Println("Error creating XML file:", err)
		return
	}
	_, err = xmlFile.Write([]byte(xml.Header)) // write xml header
	if err != nil {
		fmt.Println("Error writing to XML file:", err)
		return
	}

	encoder := xml.NewEncoder(xmlFile)
	encoder.Indent("", "\t")

	// write gpx
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
