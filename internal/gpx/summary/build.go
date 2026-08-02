package summary

import (
	"fmt"
	"html/template"
	"io"
	"slices"
	"strconv"
	"strings"
)

func buildHeaders(opt Option) []string {
	headers := []string{
		"Id",
		setValueIf(opt.WithName, "Name"),
		setValueIf(!opt.WithName, "From"),
		setValueIf(opt.WithDestination, "To"),
		setValueIf(opt.WithPoints, "N pts"),
		setValueIf(opt.WithDistance, "Distance(km)"),
		setValueIf(opt.WithAscentDescent, "Ascent(m)"),
		setValueIf(opt.WithAscentDescent, "Descent(m)"),
		setValueIf(opt.WithDistanceEffort, "DistanceEffort(km)"),
		setValueIf(opt.WithDuration, "Duration"),
	}
	headers = slices.DeleteFunc(
		headers,
		func(s string) bool {
			return s == ""
		})
	return headers
}

func buildValues(s GeoStatistic, opt Option) []string {
	values := []string{
		setValueIf(opt.WithName, fmt.Sprintf("%v", s.Name)),
		setValueIf(!opt.WithName, fmt.Sprintf("%v", s.From)),
		setValueIf(opt.WithDestination, fmt.Sprintf("%v", s.To)),
		setValueIf(opt.WithPoints, fmt.Sprintf("%v", s.N)),
		setValueIf(opt.WithDistance, fmt.Sprintf("%.0f", s.Distance)),
		setValueIf(opt.WithAscentDescent, fmt.Sprintf("%.0f", s.TotalAscent)),
		setValueIf(opt.WithAscentDescent, fmt.Sprintf("%.0f", s.TotalDescent)),
		setValueIf(opt.WithDistanceEffort, fmt.Sprintf("%.0f", s.DistanceEffort)),
		setValueIf(opt.WithDuration, fmt.Sprintf("%vh%v", s.DurationHour, s.DurationMin)),
	}

	// remove unset value
	values = slices.DeleteFunc(values, func(s string) bool { return s == "" })
	return values
}

func (s GeoStatistic) WriteHeader(w io.Writer, kind string, opt Option) error {
	switch kind {
	case "csv":
		return s.writeHeaderCsv(w, opt)
	case "html":
		return s.writeHeaderHTML(w, opt)
	case "string":
	default:
		fmt.Errorf("kind must be one of values: string, csv or html. But is '%v'", kind)
	}
	return nil
}

func (s GeoStatistic) WriteFooter(w io.Writer, kind string) error {
	switch kind {
	case "html":
		return s.writeFooterHTML(w)
	case "csv", "string":
	default:
		fmt.Errorf("kind must be one of values: string, csv or html. But is '%v'", kind)
	}
	return nil
}

func (s GeoStatistic) WriteValues(w io.Writer, i *int, kind string, opts ...Option) error {
	opt := mergeOptions(opts)

	switch kind {
	case "html":
		return s.writeValuesHTML(w, i, opt)
	case "csv":
		return s.writeValuesCsv(w, i, opt)
	case "string":
		return s.writeValuesString(w, i, opt)
	default:
		fmt.Errorf("kind must be one of values: string, csv or html. But is '%v'", kind)
	}
	return nil
}

func CsvHeader(opt Option) string {
	headers := buildHeaders(opt)
	return strings.Join(headers, ",") + "\n"
}

func (s GeoStatistic) writeHeaderCsv(w io.Writer, opt Option) (err error) {
	_, err = w.Write([]byte(CsvHeader(opt)))
	return
}
func (s GeoStatistic) writeValuesCsv(w io.Writer, i *int, opt Option) (err error) {
	var id string
	if i != nil {
		id = strconv.Itoa(*i + 1)
	}
	_, err = w.Write(
		[]byte(fmt.Sprintf("%v,%v", id, s.ToCsv(opt))))
	return
}

// func (s GeoStatistic) writeFooterCsv(w io.Writer) (err error) {}

func loadHTMLTemplate() (*template.Template, error) {
	tplPath := "internal/gpx/summary/"
	// NOTE: table*.html can be merge. Requires setting th/td as argument
	return template.ParseFiles([]string{
		tplPath + "table_header_rows.html",
		tplPath + "template.html",
		tplPath + "table_body_rows.html",
	}...)
}

func (s GeoStatistic) writeHeaderHTML(w io.Writer, opt Option) (err error) {
	t, err := loadHTMLTemplate()
	if err != nil {
		return err
	}
	headers := buildHeaders(opt)

	return t.ExecuteTemplate(
		w, "HTMLTop",
		struct{ TableHeaderRows []string }{headers})

}
func (s GeoStatistic) writeValuesHTML(w io.Writer, i *int, opt Option) (err error) {
	t, err := loadHTMLTemplate()
	if err != nil {
		return err
	}

	var id string
	if i != nil {
		id = strconv.Itoa(*i + 1)
	}
	values := slices.Concat(
		[]string{id}, // days index
		buildValues(s, opt))

	if err = t.ExecuteTemplate(
		w, "TableBodyRows",
		struct{ TableBodyRows []string }{values}); err != nil {
		return
	}

	return nil
}
func (s GeoStatistic) writeFooterHTML(w io.Writer) (err error) {
	t, err := loadHTMLTemplate()
	if err != nil {
		return err
	}
	return t.ExecuteTemplate(w, "HTMLBottom", nil)
}

// func (s GeoStatistic) writeHeaderString(w io.Writer) (err error) {}
func (s GeoStatistic) writeValuesString(w io.Writer, i *int, opt Option) (err error) {
	var id string
	if i != nil {
		id = fmt.Sprintf("[%v]", *i)
	}

	_, err = w.Write([]byte(fmt.Sprintf("%v %v\n", id, s.ToString(opt))))
	return
}

// func (s GeoStatistic) writeFooterString(w io.Writer) (err error) {}
