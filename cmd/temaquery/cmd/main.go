package cmd

import (
	"github.com/MathWebSearch/mwsapi/connection"
	"github.com/MathWebSearch/mwsapi/engine/temaengine"
	"github.com/MathWebSearch/mwsapi/query"
)

// Main represents the main interface of the temaquery command
func Main(a *Args) (res interface{}, err error) {
	// make a new connection
	c, err := connection.NewTemaConnection(a.MWSPort, a.MWSHost, a.ElasticPort, a.ElasticHost)
	if err != nil {
		return
	}

	// connect
	err = connection.Connect(c)
	if err != nil {
		return
	}
	defer c.Close()

	// make a query
	q := &query.Query{
		Expressions: a.Expressions,
		Text:        a.Text,
	}

	// run the count (if requested)
	if a.Count {
		return temaengine.Count(c, q)
	}

	{
		res, err := temaengine.Run(c, q, a.From, a.Size)

		// normalize if requested
		if err == nil && a.Normalize {
			res.Normalize()
		}

		return res, err
	}
}
