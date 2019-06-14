package main

import (
	"os"

	"github.com/MathWebSearch/mwsapi/cmd/mwsapid/cmd"
)

func main() {
	// parse and validate arguments
	a := cmd.ParseArgs(os.Args)
	if !a.Validate() {
		os.Exit(1)
		return
	}

	// run the code and exit
	cmd.Main(a)
}
