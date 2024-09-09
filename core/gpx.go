package core

import (
	"encoding/xml"
	"fmt"
	"os"
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

func (gpx *Gpx) Info(gpxFilename string, vitessePlat float64, detail bool, ascii_format bool) {
	(*gpx).ParseFile(gpxFilename)
	(*gpx).SetVitesse(vitessePlat)

	for i, trk := range gpx.Trk {
		summary := trk.CalcAll((*gpx).Extensions.Vitesse, detail)
		for _, s := range summary {

			fmt.Printf("[%v] ", i)
			s.Print()
		}
		// trk.PrintInfo(ascii_format)

		fmt.Println()
	}

}
