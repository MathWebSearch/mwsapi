package cmd

import (
	"github.com/MathWebSearch/mwsapi/connection"
	"github.com/MathWebSearch/mwsapi/engine/mwsengine"
	"github.com/MathWebSearch/mwsapi/query"
)

// Main represents the main interface of the mwsquery command
func Main(a *Args) (res interface{}, err error) {
	// make a new connection
	c, err := connection.NewMWSConnection(a.MWSPort, a.MWSHost)
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
	q := &query.MWSQuery{
		Expressions: a.Expressions,
		MwsIdsOnly:  a.MWSIdsOnly,
	}

	// run the count (if requested)
	if a.Count {
		return mwsengine.Count(c, q)
	}

	{
		res, err := mwsengine.Run(c, q, a.From, a.Size)

		// normalize if requested
		if err == nil && a.Normalize {
			res.Normalize()
		}

		return res, err
	}
}
