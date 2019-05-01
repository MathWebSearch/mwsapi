package main

import (
	"os"

	"github.com/MathWebSearch/mwsapi/cmd/temaquery/args"
	"github.com/MathWebSearch/mwsapi/connection"
	"github.com/MathWebSearch/mwsapi/engine/temaengine"
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
	c, err := connection.NewTemaConnection(a.MWSPort, a.MWSHost, a.ElasticPort, a.ElasticHost)
	if err != nil {
		panic(err)
	}

	// connect
	err = connection.Connect(c)
	if err != nil {
		panic(err)
	}

	// make a query
	q := &query.Query{
		Expressions: a.Expressions,
		Text:        a.Text,
	}

	// run
	var res interface{}

	if !a.Count {
		r, e := temaengine.Run(c, q, a.From, a.Size)

		// normalize if requested
		if e == nil && a.Normalize {
			r.Normalize()
		}

		// and store the results
		res = r
		err = e
	} else {
		res, err = temaengine.Count(c, q)
	}

	// and output
	utils.OutputJSONOrErr(&res, err)
}
