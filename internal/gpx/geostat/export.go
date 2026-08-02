package geostat

import (
	"fmt"
	"slices"
	"strings"

	sym "github.com/jple/text-symbol"
)

type Option struct {
	WithDestination    bool // Whether print destination name (To)
	WithPoints         bool // Whether print number of points (N)
	WithDistanceEffort bool // Whether print `DistanceEffort`
	WithDistance       bool // Whether print `Distance`
	WithDuration       bool // Whether print duration
	WithAscentDescent  bool // Whether print Ascent and Descent

	// Specific to string
	Indent int // number of space to indent string

	// Specific to csv
	WithName bool
	WithFrom bool

	// Specific to csv
	AddHeader bool
}

var defOpt Option = Option{
	WithDestination:    true,
	WithDistanceEffort: true,
	WithDistance:       true,
	WithDuration:       true,
	WithAscentDescent:  true,
}

func mergeOptions(opts []Option) Option {
	opt := defOpt

	if len(opts) > 0 {
		opt = Option{
			WithDestination:    slices.ContainsFunc(opts, func(opt Option) bool { return opt.WithDestination }),
			WithPoints:         slices.ContainsFunc(opts, func(opt Option) bool { return opt.WithPoints }),
			WithDistanceEffort: slices.ContainsFunc(opts, func(opt Option) bool { return opt.WithDistanceEffort }),
			WithDistance:       slices.ContainsFunc(opts, func(opt Option) bool { return opt.WithDistance }),
			WithDuration:       slices.ContainsFunc(opts, func(opt Option) bool { return opt.WithDuration }),
			WithAscentDescent:  slices.ContainsFunc(opts, func(opt Option) bool { return opt.WithAscentDescent }),

			// Specific to string
			Indent: opts[len(opts)-1].Indent,

			// Specific to csv
			WithName:  slices.ContainsFunc(opts, func(opt Option) bool { return opt.WithName }),
			WithFrom:  slices.ContainsFunc(opts, func(opt Option) bool { return opt.WithFrom }),
			AddHeader: slices.ContainsFunc(opts, func(opt Option) bool { return opt.AddHeader }),
		}
	}
	return opt
}

func setValueIf(cond bool, val string) string {
	var s string
	if cond {
		s = val
	}
	return s
}

func (s GeoStatistic) ToString(opts ...Option) string {
	opt := mergeOptions(opts)

	var str string
	str += strings.Repeat(" ", opt.Indent)
	str += sym.Green(s.From)
	str += setValueIf(opt.WithDestination, fmt.Sprintf(" --> %v", sym.Green(s.To)))

	stats := []string{
		setValueIf(opt.WithPoints, fmt.Sprintf("%.0f pts", s.N)),
		setValueIf(opt.WithDistance, fmt.Sprintf("%v %.0fkm", sym.ArrowIconLeftRight(), s.Distance)),
		setValueIf(opt.WithDistanceEffort, fmt.Sprintf("%v %.0fkm_e", sym.ArrowWaveRight(), s.DistanceEffort)),
		setValueIf(opt.WithAscentDescent, fmt.Sprintf("%v+%.0fm/%.0fm", sym.UpAndDown(), s.TotalAscent, s.TotalDescent)),
		setValueIf(opt.WithDuration, fmt.Sprintf("%v %vh%02d", sym.StopWatch(), s.DurationHour, s.DurationMin)),
	}
	stats = slices.DeleteFunc(stats, func(s string) bool { return s == "" })
	if len(stats) > 0 {
		str += fmt.Sprintf("(%v)", strings.Join(stats, ", "))
	}

	return str
}

func (s GeoStatistic) ToCsv(opts ...Option) string {
	opt := mergeOptions(opts)
	return strings.Join(buildValues(s, opt), ",") + "\n"
}
