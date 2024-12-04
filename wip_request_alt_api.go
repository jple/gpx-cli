//go:build exclude

package main

import (
	"fmt"
	"io"
	"net/http"
	"net/url"

	. "github.com/jple/gpx-cli/core"
)

type QueryParam struct {
	url       string
	lon       string
	lat       string
	resource  string
	delimiter string
	zonly     string
}

func GetElevation(param QueryParam) []byte {
	client := &http.Client{}

	if param.delimiter == "" {
		param.delimiter = ","
	}
	if param.zonly == "" {
		param.zonly = "true"
	}

	// Create query
	req, _ := http.NewRequest("GET", param.url, nil)

	req.URL.RawQuery = url.Values{
		"lon":       {param.lon},
		"lat":       {param.lat},
		"resource":  {param.resource},
		"delimiter": {param.delimiter},
		"zonly":     {param.zonly},
	}.Encode()

	req.Header.Add("accept", "application/json")

	fmt.Println("request url:", req.URL)
	fmt.Println("request header:", req.Header)

	// Get results
	resp, _ := client.Do(req)
	body, _ := io.ReadAll(resp.Body)
	return body
}

const AltimetrieApiUrl = "https://data.geopf.fr/altimetrie/1.0/calcul/alti/rest/elevation.json"

func main() {
	gpx := Gpx{Filepath: "src/my.gpx"}
	gpx.ParseFile(gpx.Filepath)

	lons, lats := gpx.Trk[0].GetLonLat()
	fmt.Println(len(lons))
	fmt.Println(len(lats))

	// query := QueryParam{
	// 	url:      AltimetrieApiUrl,
	// 	lon:      strings.Join(lons, ","),
	// 	lat:      strings.Join(lats, ","),
	// 	resource: "ign_rge_alti_par_territoires",
	// }

	// body := GetElevation(query)
	// fmt.Printf("%s", body)
}
