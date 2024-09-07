package main

import (
	cmd "github.com/jple/gpx-cli/cmd"
)

func main() {
	// Usage: coucou
	// Usage: coucou cmd1
	// Usage: coucou --help
	// Usage: coucou cmd1 --help
	cmd.Execute()
}
