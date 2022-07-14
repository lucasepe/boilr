package main

import (
	"fmt"
	"os"

	"github.com/lucasepe/boilr/cmd"
)

// nolint: gochecknoglobals
var (
	Version string
	Build   string
)

func main() {
	err := cmd.Run(Version, Build)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err.Error())
		os.Exit(1)
	}
}
