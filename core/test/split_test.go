package core

import (
	"fmt"
	"os/exec"
	"testing"

	. "github.com/jple/gpx-cli/core"
)

// TODO: this function only exists because GetInfo does not return correct NPoints...
// and only for dev purpose of TestSplit
// func printNTrkpt(gpx Gpx) {
// 	for i, trk := range gpx.Trk {
// 		fmt.Printf("Trk [%v]", i)
// 		n := 0
// 		for _, trkseg := range trk.Trkseg {
// 			for range trkseg.Trkpt {
// 				n++
// 			}
// 		}
// 		fmt.Printf(" #Trkpt: *%v*\n", n)
// 	}
// }

func TestSplit(t *testing.T) {
	gpx := Gpx{}
	gpx.ParseFile("data/split.gpx")
	fmt.Println("=======before split=========")
	gpx.GetInfo(true).ToString(PrintArgs{})
	// fmt.Println("************ TEST  **********")
	// fmt.Printf("%+v\n", gpx)

	var cmd *exec.Cmd
	for i := 0; i <= 3; i++ {
		cmd = exec.Command("rm", fmt.Sprintf("out-%v.gpx", i))
		cmd.Run()
		cmd = exec.Command("rm", fmt.Sprintf("gpx-split-%v.gpx", i))
		cmd.Run()
	}
	cmd.Run()

	fmt.Println("=======after split=========")
	gpx.Save("out-0.gpx")

	gpx1 := gpx.SplitAtName("first")
	fmt.Println("Split first")
	gpx1.GetInfo(true).ToString(PrintArgs{})
	gpx.Save("gpx-split-1.gpx")
	gpx1.Save("out-1.gpx")

	gpx2 := gpx.SplitAtName("between")
	fmt.Println("Split between")
	gpx2.GetInfo(true).ToString(PrintArgs{})
	gpx.Save("gpx-split-2.gpx")
	gpx2.Save("out-2.gpx")

	gpx3 := gpx.SplitAtName("last")
	fmt.Println("Split last")
	gpx3.GetInfo(true).ToString(PrintArgs{})
	// gpx3 := gpx.Split(0, 2, 0) // BUG: gpx is actually updated...
	// fmt.Println("************ gpx-split **********")
	// fmt.Printf("%+v\n", gpx)
	// fmt.Println("************ out **********")
	// fmt.Printf("%+v\n", gpx3)

	gpx.Save("gpx-split-3.gpx")
	gpx3.Save("out-3.gpx")

	// BUG: does it come from "first" name at trk ??

	// fmt.Println(len(gpx1.Trk) == len(gpx.Trk))
	// fmt.Println(len(gpx2.Trk) == len(gpx.Trk)+1)
	// fmt.Println(len(gpx3.Trk) == len(gpx.Trk)+1)

	// fmt.Println(len(gpx2.Trk[0].Trkseg) == 1)
	// fmt.Println(len(gpx3.Trk[0].Trkseg) == len(gpx.Trk[0].Trkseg)-1)
	// fmt.Println(len(gpx3.Trk[0].Trkseg))
	// fmt.Println(len(gpx.Trk[0].Trkseg))

	// gpx.GetInfo(true).ToString(PrintArgs{})

	// input := 0
	// out := f(input)
	// wants := 1

	// if out != wants {
	// 	t.Fatalf(`f(%v) returns %v, wants %v`, input, out, wants)
	// }
}
