package ign

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/jple/gpx-cli/internal/geo"
)

const (
	apiUrl string = "https://data.geopf.fr/altimetrie/1.0/calcul/alti/rest/elevation.json"

	defaultResource  string = "ign_rge_alti_wld"
	defaultDelimiter string = ";"
)

type IgnOutput struct {
	Elevations []float64 `json:"elevations"`
}

type IgnRequest struct {
	// NOTE: reuse geo.Coord !!
	Lat       string `json:"lat"`
	Lon       string `json:"lon"`
	Resource  string `json:"resource"`
	Delimiter string `json:"delimiter"`
	Measures  string `json:"measures"`
	Zonly     string `json:"zonly"`
	Indent    string `json:"indent"`
}

type Coords []geo.Coord
type Values url.Values

func (pts Coords) getLatLonFormated(delim string) (string, string) {
	var lats, lons []string
	for _, p := range pts {
		lats = append(lats, strconv.FormatFloat(p.Lat, 'f', -1, 64))
		lons = append(lons, strconv.FormatFloat(p.Lon, 'f', -1, 64))
	}

	return strings.Join(lats, delim), strings.Join(lons, delim)

}

func (v Values) ToString() string {
	var out string
	for k := range v {
		nextParam := k + "=" + strings.Join(v[k], "")
		if out == "" {
			out = nextParam
		} else {
			out += "&" + nextParam
		}
	}
	return out
}

func Elevations(pts Coords) []float64 {
	// Prepare request
	req, _ := http.NewRequest("GET", apiUrl, nil)
	lats, lons := pts.getLatLonFormated(defaultDelimiter)
	req.URL.RawQuery = Values{
		"lat":       {lats},
		"lon":       {lons},
		"resource":  {defaultResource},
		"delimiter": {defaultDelimiter},
		"measures":  {"false"},
		"zonly":     {"true"},
	}.ToString()

	// Query
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}

	// Format
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}
	var ignResp IgnOutput
	err = json.Unmarshal(body, &ignResp)
	if err != nil {
		fmt.Println("Query failed")
		fmt.Println(string(body))
	}

	return ignResp.Elevations
}

func PostElevations(pts Coords) []float64 {
	// Prepare body request to query
	lats, lons := pts.getLatLonFormated(defaultDelimiter)
	ignReq := IgnRequest{
		Lat:       lats,
		Lon:       lons,
		Delimiter: defaultDelimiter,
		Resource:  defaultResource,
		Zonly:     "true",
		Measures:  "false",
		Indent:    "false",
	}
	j, _ := json.Marshal(ignReq)
	reqbody := bytes.NewBuffer(j)

	// Prepare request
	req, _ := http.NewRequest("POST", apiUrl, reqbody)

	// Query
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}

	// Format
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(resp.StatusCode)
	if resp.StatusCode != 200 {
		fmt.Println(string(body))
	}

	var ignResp IgnOutput
	err = json.Unmarshal(body, &ignResp)
	if err != nil {
		fmt.Println("Query failed")
		fmt.Println(string(body))
	}

	return ignResp.Elevations
}

func F() {
	pts := Coords{
		geo.Coord{Lat: 44.4, Lon: 3.2},
		geo.Coord{Lat: 43.2, Lon: 1.3},
	}
	// elevations := Elevations(pts)
	// fmt.Println(elevations)

	elevations := PostElevations(pts)
	fmt.Println(elevations)
}
