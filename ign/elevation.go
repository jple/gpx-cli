package ign

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

const (
	apiUrl string = "https://data.geopf.fr/altimetrie/1.0/calcul/alti/rest/elevation.json"

	defaultResource  string = "ign_rge_alti_wld"
	defaultDelimiter string = ";"
)

type (
	IgnOutput struct {
		Elevations []float64 `json:"elevations"`
	}

	IgnRequest struct {
		Lat       string `json:"lat"`
		Lon       string `json:"lon"`
		Resource  string `json:"resource"`
		Delimiter string `json:"delimiter"`
		Measures  bool   `json:"measures"`
		Zonly     bool   `json:"zonly"`
		Indent    bool   `json:"indent"`
	}

	Point struct {
		Lat float64 `json:"lat"`
		Lon float64 `json:"lon"`
	}

	Points []Point

	Values url.Values
)

func (pts Points) getLatLonFormated(delim string) (string, string) {
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

func GetElevations(pts Points) IgnOutput {
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
	resp, _ := client.Do(req)

	// Format
	body, _ := io.ReadAll(resp.Body)
	var ignResp IgnOutput
	json.Unmarshal(body, &ignResp)

	return ignResp
}

func F() {
	pts := Points{
		Point{Lat: 44.4, Lon: 3.2},
		Point{Lat: 43.2, Lon: 1.3},
	}
	GetElevations(pts)
}
