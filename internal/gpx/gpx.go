package gpx

import (
	"encoding/xml"
	"errors"
	"fmt"
	"os"
	"slices"
	"strings"
)

// TODO: change fmt.Print into log
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

func check(err error) {
	if err != nil {
		fmt.Println(err)
	}
}

// ------------------- I/O ------------------------

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

// ------------------- Query ------------------------

func (gpx Gpx) AllTrkpts() (trkpts Trkpts) {
	for _, trk := range gpx.Trks {
		trkpts = slices.Concat(trkpts, trk.AllTrkpts())
	}
	return
}

func (gpx Gpx) FindTrkptsId(trkId, segId, ptId int) (id int) {
	// Add all trkpts in Trk before Trk[TrkId]
	for i := 0; i < trkId; i++ {
		id += len(gpx.Trks[i].AllTrkpts())
	}
	// Add all trkpts in Seg before Seg[SegId], from Trk[trkId]
	for i := 0; i < segId; i++ {
		id += len(gpx.Trks[trkId].Trksegs[i].Trkpts)
	}
	// Add all trkpts before ptId (included), from Trk[trkId].Seg[segId]
	id += ptId

	return id
}

// NOTE: compare speed between AllTrkpts() + loop for name vs. loop trk, seg, pt for name
func (gpx Gpx) FindTrkptsIdByName(name string) (int, error) {
	trkpts := gpx.AllTrkpts()
	for id, pt := range trkpts {
		if pt.Name != nil && *pt.Name == name {
			return id, nil
		}
	}
	return -1, errors.New("Name not found")
}

func GpxTrkptsBetweenIndex(gpx Gpx, i1, i2 int) (Trkpts, error) {
	allTrkpts := gpx.AllTrkpts()
	if i1 < 0 || i2 < 0 {
		return nil, fmt.Errorf("i1 and i2 must be >= 0\n")
	}
	if i1 >= len(allTrkpts) || i2 >= len(allTrkpts) {
		return nil, fmt.Errorf("i1 and i2 must be < the number of all Trkpts\n")
	}

	var a, b int
	if i1 <= i2 {
		a = i1
		b = i2
	} else {
		a = i2
		b = i1
	}

	return allTrkpts[a : b+1], nil
}

func GpxTrkptsBetweenNames(gpx Gpx, name1, name2 string) (Trkpts, error) {
	var i1, i2 int
	var err error
	if i1, err = gpx.FindTrkptsIdByName(name1); err != nil {
		return nil, err
	}
	if i2, err = gpx.FindTrkptsIdByName(name2); err != nil {
		return nil, err
	}
	return GpxTrkptsBetweenIndex(gpx, i1, i2)
}

// ------------------- Modifications ------------------------
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

// Split Trk[trkId] containing trkptId into two trk 0:trptkId (excluded) and trkptId:end (included)
func (gpx *Gpx) Split(trkId, trksegId, trkptId int) {
	// Input validation
	if trkId < 0 || trkId >= len(gpx.Trks) {
		return
	}
	if trksegId < 0 || trksegId >= len(gpx.Trks[trkId].Trksegs) {
		return
	}
	if trkptId < 0 || trkptId >= len(gpx.Trks[trkId].Trksegs[trksegId].Trkpts) {
		return
	}

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
			out.Trksegs = slices.Concat(out.Trksegs, []Trkseg{trkseg_last}) // NOTE: using append would modify gpx, which is not wanted
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

	bef := filterBeforeTrkpt(*gpx, trkId, trksegId, trkptId, gpx.Trks[trkId].Name)
	aft := filterAfterTrkpt(*gpx, trkId, trksegId, trkptId, *gpx.Trks[trkId].Trksegs[trksegId].Trkpts[trkptId].Name)

	// Update gpx.Trks value
	newTrk := slices.Clone(gpx.Trks) // prevent unwanted update on gpx argument
	newTrk = slices.Delete(newTrk, trkId, trkId+1)
	if len(bef.Trksegs) > 0 {
		newTrk = slices.Insert(newTrk, trkId, bef)
		newTrk = slices.Insert(newTrk, trkId+1, aft)
	} else {
		newTrk = slices.Insert(newTrk, trkId, aft)
	}
	gpx.Trks = newTrk

}

// SplitAtName returns a Gpx with a trk split at the trkpt name : 0:trkptName (excluded) and trkptName:end (included)
func (gpx *Gpx) SplitAtName(name string) {
	found := false

	for i, trk := range gpx.Trks {
		for j, trkseg := range trk.Trksegs {
			for k, trkpt := range trkseg.Trkpts {
				if trkpt.Name != nil && *trkpt.Name == name {
					found = true

					// TODO: this print is a pb in tui module
					// fmt.Printf("Split at trk %v, trkseg %v, trkpt %v\n", i, j, k)

					gpx.Split(i, j, k)
					return

				}
			}
		}
	}

	if !found {
		// TODO: this print is a pb in tui module
		// fmt.Printf("Name '%v' not found in gpx\n", name)
	}
	return
}
