package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"text/tabwriter"

	. "github.com/jple/gpx-cli/core"
)

func prettyStruct(in any) string {
	j, err := json.MarshalIndent(in, "", "  ")
	// j, err := json.Marshal(in) // to return as data-raw
	if err != nil {
		log.Fatalf(err.Error())
	}
	return string(j)
}

func test() {
	gpx := Gpx{}
	gpx.ParseFile("src/my.gpx")
	gpx.SetVitesse(4.5)

	var printArgs PrintArgs = PrintArgs{AsciiFormat: true}
	trkid, _ := strconv.ParseInt(os.Args[1], 0, 0)

	w := tabwriter.NewWriter(os.Stdout, 1, 1, 1, ' ', tabwriter.Debug)
	if trkid == -1 {
		printArgs.Silent = true
		infos := gpx.
			GetInfo(true).
			ToString(printArgs)
		fmt.Fprintln(w, infos)
		// fmt.Println(str)
		// fmt.Printf("%+v\n", gpx.GetInfo(true))
		// fmt.Println(prettyStruct(gpx.GetInfo(true)))
	} else {
		printArgs.PrintFrom = true
		// fmt.Println()
		// fmt.Println()
		infos := gpx.
			Trk[trkid].
			GetInfo(gpx.Extensions.Vitesse, true).
			ToString(printArgs)
			// fmt.Println(str)
		fmt.Fprintln(w, infos)
		// fmt.Println(prettyStruct(gpx.Trk[trkid].GetInfo(gpx.Extensions.Vitesse, true)))
	}
	w.Flush()
}

func main() {
	// cmd.Execute()
	test()
	// sym.ShowUnicode()
}
