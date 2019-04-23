package main

import (
	"os"

	"github.com/MathWebSearch/mwsapi/cmd/mwsquery/args"
	"github.com/MathWebSearch/mwsapi/mws"
	"github.com/MathWebSearch/mwsapi/utils"
)

func main() {
	// parse and validate arguments
	a := args.ParseArgs(os.Args)
	if !a.Validate() {
		os.Exit(1)
		return
	}

	connection := mws.NewConnection(a.MWSHost, a.MWSPort)
	query := &mws.Query{
		Expressions: a.Expressions,
		MwsIdsOnly:  a.MWSIdsOnly,
	}

	var res interface{}
	var err error

	if a.Count {
		res, err = mws.CountQuery(connection, query)
	} else {
		res, err = mws.RunQuery(connection, query, a.From, a.Size)
	}

	// and output
	utils.OutputJSONOrErr(&res, err)
}
