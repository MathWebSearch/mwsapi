package main

import (
	"os"

	"github.com/MathWebSearch/mwsapi/result"

	"github.com/MathWebSearch/mwsapi/cmd/elasticquery/args"
	"github.com/MathWebSearch/mwsapi/connection"
	"github.com/MathWebSearch/mwsapi/engine/elasticengine"
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
	c, err := connection.NewTemaConnection(a.ElasticPort, a.ElasticHost)
	if err != nil {
		panic(err)
	}

	// connect
	err = connection.Connect(c)
	if err != nil {
		panic(err)
	}

	// make a query
	q := &query.ElasticQuery{
		MathWebSearchIDs: a.IDs,
		Text:             a.Text,
	}

	// keep the results
	var res interface{}

	if !a.Count {

		// run either the entire thing or the document query
		var r *result.Result
		var e error
		if !a.DocumentPhaseOnly {
			r, e = elasticengine.Run(c, q, a.From, a.Size)
		} else {
			r, e = elasticengine.RunDocument(c, q, a.From, a.Size)
		}

		// normalize if requested
		if e == nil && a.Normalize {
			r.Normalize()
		}

		// and store the results
		res = r
		err = e

		// count
	} else {
		res, err = elasticengine.Count(c, q)
	}

	// and output
	utils.OutputJSONOrErr(&res, err)
}
