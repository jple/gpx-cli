module github.com/jple/gpx-cli

go 1.23.4

replace github.com/jple/gpx-cli/cmd => ./cmd

replace github.com/jple/gpx-cli/core => ./core

replace github.com/jple/gpx-cli/ign => ./ign

replace github.com/jple/gpx-cli/tui => ./tui

require github.com/jple/gpx-cli/core v0.0.0-20250406071217-76fdcb001c54

require (
	codeberg.org/go-fonts/liberation v0.5.0 // indirect
	codeberg.org/go-latex/latex v0.1.0 // indirect
	codeberg.org/go-pdf/fpdf v0.10.0 // indirect
	git.sr.ht/~sbinet/gg v0.6.0 // indirect
	github.com/ajstarks/svgo v0.0.0-20211024235047-1546f124cd8b // indirect
	github.com/campoy/embedmd v1.0.0 // indirect
	github.com/golang/freetype v0.0.0-20170609003504-e2365dfdc4a0 // indirect
	github.com/jple/text_symbol v0.0.0-00010101000000-000000000000 // indirect
	github.com/pmezard/go-difflib v1.0.1-0.20181226105442-5d4384ee4fb2 // indirect
	golang.org/x/image v0.25.0 // indirect
	golang.org/x/text v0.23.0 // indirect
	gonum.org/v1/plot v0.16.0 // indirect
)

replace github.com/jple/text_symbol => ../text_symbol
