package cmd

import (
	"github.com/MathWebSearch/mwsapi/connection"
	"github.com/MathWebSearch/mwsapi/engine/temaengine"
	"github.com/MathWebSearch/mwsapi/query"
	"github.com/pkg/errors"
)

// Main represents the main interface of the temaquery command
func Main(a *Args) (res interface{}, err error) {
	// make a new connection
	c, err := connection.NewTemaConnection(a.MWSPort, a.MWSHost, a.ElasticPort, a.ElasticHost)
	if err != nil {
		err = errors.Wrap(err, "connection.NewTemaConnection failed")
		return
	}

	// connect
	err = connection.Connect(c)
	if err != nil {
		err = errors.Wrap(err, "connection.Connect failed")
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
		count, err := temaengine.Count(c, q)
		err = errors.Wrap(err, "temaengine.Count failed")
		return count, err
	}

	{
		res, err := temaengine.Run(c, q, a.From, a.Size)
		err = errors.Wrap(err, "temaengine.Run failed")

		// normalize if requested
		if err == nil && a.Normalize && res != nil {
			res.Normalize()
		}

		return res, err
	}
}
