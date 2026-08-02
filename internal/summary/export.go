package summary

import (
	"fmt"
	"strconv"
)

func (gpxSummary GpxSummary) ToCsv() string {
	var str string

	// Add csv header
	str += fmt.Sprintln("Id,From,To,N points,Distance(km),Elevation+(m),Elevation-(m),Distance effort (km),Duration")
	for i, trkSummary := range gpxSummary.TrkSummaries {

		// Add trk summary
		str += strconv.Itoa(i+1) + "," + trkSummary.TrkName + ","
		str += fmt.Sprintf("%v,%v,%.0f,%.0f,%.0f,%.0f,%vh%v\n",
			"", // To value set to empty. "From" value is already set earlier
			trkSummary.AllSections.Len(),
			trkSummary.AllSections.Distance,
			trkSummary.AllSections.TotalAscent, trkSummary.AllSections.TotalDescent,
			trkSummary.AllSections.DistanceEffort,
			trkSummary.AllSections.DurationHour, trkSummary.AllSections.DurationMin)

		// Add section summary
		if len(trkSummary.PerSection) > 1 {
			for _, sectionSummary := range trkSummary.PerSection {
				str += fmt.Sprintf(",%v,%v,%v,%.0f,%.0f,%.0f,%.0f,%vh%v\n",
					sectionSummary.From, sectionSummary.To,
					sectionSummary.Len(),
					sectionSummary.Distance,
					sectionSummary.TotalAscent, sectionSummary.TotalDescent,
					sectionSummary.DistanceEffort,
					sectionSummary.DurationHour, sectionSummary.DurationMin)
			}
		}
	}
	return str
}

func (gpxSummary GpxSummary) ToHTML() string {
	var str string

	// Add csv header
	str += `
<style type="text/css">
table, th, td  {
    border: 2px solid lightgrey;
    border-collapse: collapse;
}
thead th {
    border: 2px double black;
    font-weight: bold;
}
tfoot td {
    border: 2px double black;
    font-weight: bold;
}
</style>

<table>
    <thead>
        <tr>
            <th>Id</th>
            <th>From</th>
            <th>To</th>
            <th>N Points</th>
            <th>Distance (km)</th>
            <th>Elevation + (m)</th>
            <th>Elevation - (m)</th>
            <th>Distance effort (km)</th>
            <th>Duration</th>
        </tr>
    </thead>
<tbody>
`
	var totalDistance, totalDistanceEffort,
		totalAscent, totalDescent float64

	for i, trkSummary := range gpxSummary.TrkSummaries {
		totalDistance += trkSummary.AllSections.Distance
		totalAscent += trkSummary.AllSections.TotalAscent
		totalDescent += trkSummary.AllSections.TotalDescent
		totalDistanceEffort += trkSummary.AllSections.DistanceEffort

		// Add trk summary
		str += fmt.Sprintf(`
			<tr>
				<td>%v</td>
				<td>%v</td>
				<td>%v</td>
				<td>%v</td>
				<td>%.0f</td>
				<td>%.0f</td>
				<td>%.0f</td>
				<td>%.0f</td>
				<td>%vh%v</td>
			</tr>`,
			strconv.Itoa(i+1), trkSummary.TrkName, "", // "To" value set to empty
			trkSummary.AllSections.Len(),
			trkSummary.AllSections.Distance,
			trkSummary.AllSections.TotalAscent, trkSummary.AllSections.TotalDescent,
			trkSummary.AllSections.DistanceEffort,
			trkSummary.AllSections.DurationHour, trkSummary.AllSections.DurationMin)

		// Add section summary
		if len(trkSummary.PerSection) > 1 {
			for _, sectionSummary := range trkSummary.PerSection {
				str += fmt.Sprintf(`<tr>
					<td>%v</td>
					<td>%v</td>
					<td>%v</td>
					<td>%v</td>
					<td>%.0f</td>
					<td>%.0f</td>
					<td>%.0f</td>
					<td>%.0f</td>
					<td>%vh%v</td>
				</tr>`,
					"", sectionSummary.From, sectionSummary.To,
					sectionSummary.Len(),
					sectionSummary.Distance,
					sectionSummary.TotalAscent, sectionSummary.TotalDescent,
					sectionSummary.DistanceEffort,
					sectionSummary.DurationHour, sectionSummary.DurationMin)
			}
		}
	}
	str += fmt.Sprintf(`</tbody>
		<tfoot><tr>
			<td>TOTAL</td>
			<td></td>
			<td></td>
			<td></td>
			<td>%.0f</td>
			<td>%.0f</td>
			<td>%.0f</td>
			<td>%.0f</td>
			<td></td>
		</tr></tfoot></table>`,
		totalDistance, totalAscent, totalDescent, totalDistanceEffort)

	return str
}
