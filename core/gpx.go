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
	(*gpx).Extensions.Vitesse = v
}

func (p_gpx *Gpx) Info(vitessePlat float64, detail bool, ascii_format bool) {
	gpx := *p_gpx
	gpx.ParseFile(gpx.Filepath)
	gpx.SetVitesse(vitessePlat)

	for i, trk := range gpx.Trk {
		summary := trk.CalcAll(gpx.Extensions.Vitesse, detail)
		for _, s := range summary {

			fmt.Printf("[%v] ", i)
			s.Print()
		}
		// trk.PrintInfo(ascii_format)

		fmt.Println()
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
