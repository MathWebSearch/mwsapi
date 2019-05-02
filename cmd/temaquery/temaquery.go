package main

import (
	"os"

	"github.com/MathWebSearch/mwsapi/cmd/temaquery/cmd"
	"github.com/MathWebSearch/mwsapi/utils"
)

func main() {
	// parse and validate arguments
	a := cmd.ParseArgs(os.Args)
	if !a.Validate() {
		os.Exit(1)
		return
	}

	// run the code and exit
	res, err := cmd.Main(a)
	utils.OutputJSONOrErr(res, err)
}
