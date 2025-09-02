package core

import (
	"encoding/xml"
	"errors"
	"fmt"
	"os"
	"slices"
	"strings"
)

type GpxMetadata struct {
	Name string `xml:"name,omitempty"`
	Desc *struct {
		Inner string `xml:",innerxml"`
	} `xml:"desc,omitempty"`
	Author string `xml:"author,omitempty"`
	Email  string `xml:"email,omitempty"`
	Link   []struct {
		Inner string `xml:",innerxml"`
	} `xml:"link,omitempty"`
	Url      string `xml:"url,omitempty"`
	Urlname  string `xml:"urlname,omitempty"`
	Time     string `xml:"time,omitempty"`
	Keywords []struct {
		Inner string `xml:",innerxml"`
	} `xml:"keywords,omitempty"`
	Bounds string `xml:"bounds,omitempty"`
}

type Comment struct {
	Content string `xml:",comment"`
}

type Gpx struct {
	XMLName xml.Name `xml:"gpx"`
	// XMLName xml.Name `xml:"http://www.topografix.com/GPX/1/1 gpx"`

	// --- gpx attributes ---
	// Attrs   []xmlAttr `xml:",any,attr"` // known issues: xmlns:_xmlns created multiple times
	Xmlns    string `xml:"xmlns,attr"`
	Version  string `xml:"version,attr"`
	Creator  string `xml:"creator,attr"`
	XmlnsXsi string `xml:"xmlns xsi,attr,omitempty"`
	XmlnsGh  string `xml:"xmlns gh,attr,omitempty"` // Graphhopper
	// should be xsi:schemaLocation, but not working...
	XsiSchemaLocation string `xml:"schemaLocation,attr,omitempty"`

	// --- gpx metadata ---
	Metadata *struct {
		GpxMetadata
		Comment // NOTE: comment is only permit here...
	} `xml:"metadata,omitempty"`
	GpxMetadata

	// --- gpx items ---
	Extensions *struct {
		Inner string `xml:",innerxml"`
	} `xml:"extensions,omitempty"`

	Trks []Trk `xml:"trk,omitempty"`
	Wpts []Wpt `xml:"wpt,omitempty"`
}

func check(err error) {
	if err != nil {
		fmt.Println(err)
	}
}

// Parse reads file when "in" is string (representing filepath)
// or get content when "in" is []byte (representing gpx content)
func (gpx *Gpx) Parse(in any) *Gpx {
	var data []byte

	switch d := in.(type) {
	// in is gpx filepath
	case string:
		data, _ = os.ReadFile(d)

	// in is gpx content
	case []byte:
		data = d
	}

	if err := xml.Unmarshal(data, &gpx); err != nil {
		if err.Error() != "EOF" {
			fmt.Println(err)
		}
	}

	// TODO: create receiver directly where it has to be cleaned
	// 		 then create interface to clean all
	// Cleaning struct
	if gpx.Extensions != nil && strings.TrimSpace(gpx.Extensions.Inner) == "" {
		gpx.Extensions = nil
	}

	return gpx
}

func (gpx *Gpx) AddWpt(wpt Wpt) Gpx {
	gpx.Wpts = append(gpx.Wpts, wpt)
	return *gpx
}

func (gpx *Gpx) AddColor() *Gpx {
	// TODO: create a single struct containing color, and dash*
	colors := []string{"8e44ad", "ff5733"}
	dasharray := []int{10, 10}
	dashoffset := []int{0, 5}

	for i, _ := range gpx.Trks {
		newLineColor := ExtensionsLine{
			Attrs: []xml.Attr{
				xml.Attr{
					xml.Name{"", "xmlns"},
					"http://www.topografix.com/GPX/gpx_style/0/2",
				}},
			Color:      colors[i%len(colors)],
			Dasharray:  &dasharray[i%len(dasharray)],
			Dashoffset: &dashoffset[i%len(dashoffset)],
		}
		// newLineColor.Color = colors[i%len(colors)]

		// TODO: improvement: create Trk.AddLineColor
		if gpx.Trks[i].Extensions == nil {
			fmt.Println(newLineColor)
			gpx.Trks[i].Extensions = &ExtensionsTrk{Line: &newLineColor}
		} else {
			gpx.Trks[i].Extensions.Line = &newLineColor
		}
	}
	return gpx
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

// NOTE: rename to GetInfoWhole ?
func (gpx Gpx) GetInfo(vitessePlat float64) GpxSummary {
	gpxSummary := GpxSummary{VitessePlat: vitessePlat}
	for i, trk := range gpx.Trks {
		trkSummary := trk.GetInfo(i, vitessePlat)
		gpxSummary.Trks = append(gpxSummary.Trks, struct {
			Name string
			TrkSummary
		}{trk.Name, trkSummary})
	}
	return gpxSummary
}

// ======= Get info between name ====
func (gpx Gpx) GetTrkpts() (trkpts Trkpts) {
	for _, trk := range gpx.Trks {
		trkpts = slices.Concat(trkpts, trk.GetTrkpts())
	}
	return
}

func (gpx Gpx) FindTrkptsId(trkId, segId, ptId int) (id int) {
	// Add all trkpts in Trk before Trk[TrkId]
	for i := 0; i < trkId; i++ {
		id += len(gpx.Trks[i].GetTrkpts())
	}
	// Add all trkpts in Seg before Seg[SegId], from Trk[trkId]
	for i := 0; i < segId; i++ {
		id += len(gpx.Trks[trkId].Trksegs[i].Trkpts)
	}
	// Add all trkpts before ptId (included), from Trk[trkId].Seg[segId]
	id += ptId

	return id
}

// NOTE: compare speed between gettrkpts() + loop for name vs. loop trk, seg, pt for name
func (gpx Gpx) FindTrkptsIdByName(name string) (int, error) {
	trkpts := gpx.GetTrkpts()
	for id, pt := range trkpts {
		if pt.Name != nil && *pt.Name == name {
			return id, nil
		}
	}
	return -1, errors.New("Name not found")
}

func (gpx Gpx) GetInfoBetweenTrkptsId(i1, i2 int, vitessePlat float64) TrkptsSummary {
	trkpts := gpx.GetTrkpts()
	summary := trkpts[i1 : i2+1].GetSummary(vitessePlat)

	if trkpts[i1].Name != nil {
		summary.SetFrom(*trkpts[i1].Name)
	}
	if trkpts[i2].Name != nil {
		summary.SetTo(*trkpts[i2].Name)
	}

	return summary
}

func (gpx Gpx) GetInfoBetweenName(name1, name2 string, vitessePlat float64) TrkptsSummary {
	i1, err := gpx.FindTrkptsIdByName(name1)
	check(err)
	i2, err := gpx.FindTrkptsIdByName(name2)
	check(err)
	return gpx.GetInfoBetweenTrkptsId(i1, i2, vitessePlat)
}

func (p_gpx *Gpx) Reverse() Gpx {
	gpx := *p_gpx

	slices.Reverse(gpx.Trks)
	for _, trk := range gpx.Trks {
		trk.Reverse()
	}

	return gpx
}

// Merge Trk[trkId2] into Trk[trkId1]
func (gpx *Gpx) Merge(trkId1, trkId2 int) Gpx {
	gpx.Trks[trkId1].Trksegs = slices.Concat(gpx.Trks[trkId1].Trksegs, gpx.Trks[trkId2].Trksegs)
	gpx.Trks = slices.Delete(gpx.Trks, trkId2, trkId2+1)
	return *gpx
}

// ======= WIP =========
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
		fmt.Printf("Name '%v' not found in gpx\n", name)
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
