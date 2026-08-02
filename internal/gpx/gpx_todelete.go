package gpx

import (
	"fmt"

	"github.com/jple/gpx-cli/internal/geo"
)

// TODO: To be removed after update tui/model.go
// ================= Ls and Print ===========================
type Trkname struct {
	Id        int     // Trk id in gpx
	Name      string  // Trk name
	TrksegId  *int    // Trksed id for trkpt
	TrkptId   *int    // Trkpt id for trkpt
	TrkptName *string // Trkpts name for trkpt
	*geo.Coord
	// Lat, Lon, Elevation *float64
}
type TrkNames []Trkname

func (tn Trkname) IsTrkpt() bool {
	if tn.TrkptName != nil {
		return true
	}
	return false
}
func (tn Trkname) IsTrk() bool {
	return !tn.IsTrkpt()
}

func (gpx Gpx) Ls(all bool) TrkNames {
	var out TrkNames
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
								Coord: &geo.Coord{
									Lat:       trkpt.Lat,
									Lon:       trkpt.Lon,
									Elevation: trkpt.Elevation,
								},
							})
					}
				}
			}
		}
	}

	return out
}

func (tnList TrkNames) Print(all bool, ascii_format ...bool) {

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
					trkname.Coord.Lat, trkname.Coord.Lon, trkname.Coord.Elevation)
			}
		}
	}
}
