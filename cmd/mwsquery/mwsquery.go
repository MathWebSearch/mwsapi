package main

import (
	"os"

	"github.com/MathWebSearch/mwsapi/connection"
	"github.com/MathWebSearch/mwsapi/engine/mwsengine"

	"github.com/MathWebSearch/mwsapi/cmd/mwsquery/args"
	"github.com/MathWebSearch/mwsapi/query"
	"github.com/MathWebSearch/mwsapi/utils"
)

func main() {
	// parse and validate arguments
	a := args.ParseArgs(os.Args)
	if !a.Validate() {
		os.Exit(1)
		return
	}

	// make a new connection
	c, err := connection.NewMWSConnection(a.MWSPort, a.MWSHost)
	if err != nil {
		panic(err)
	}

	// connect
	err = connection.Connect(c)
	if err != nil {
		panic(err)
	}

	// make a query
	q := &query.MWSQuery{
		Expressions: a.Expressions,
		MwsIdsOnly:  a.MWSIdsOnly,
	}

	// run
	var res interface{}

	if !a.Count {
		r, e := mwsengine.Run(c, q, a.From, a.Size)

		// normalize if requested
		if e == nil && a.Normalize {
			r.Normalize()
		}

		// and store the results
		res = r
		err = e
	} else {
		res, err = mwsengine.Count(c, q)
	}

	// and output
	utils.OutputJSONOrErr(&res, err)
}
