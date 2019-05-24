package cmd

import (
	"github.com/MathWebSearch/mwsapi/connection"
	"github.com/MathWebSearch/mwsapi/engine/mwsengine"
	"github.com/MathWebSearch/mwsapi/query"
	"github.com/pkg/errors"
)

// Main represents the main interface of the mwsquery command
func Main(a *Args) (res interface{}, err error) {
	// make a new connection
	c, err := connection.NewMWSConnection(a.MWSPort, a.MWSHost)
	if err != nil {
		err = errors.Wrap(err, "connection.NewMWSConnection failed")
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
	q := &query.MWSQuery{
		Expressions: a.Expressions,
		MwsIdsOnly:  a.MWSIdsOnly,
	}

	// run the count (if requested)
	if a.Count {
		count, err := mwsengine.Count(c, q)
		err = errors.Wrap(err, "mwsengine.Count failed")
		return count, err
	}

	{
		res, err := mwsengine.Run(c, q, a.From, a.Size)
		err = errors.Wrap(err, "mwsengine.Run failed")

		// normalize if requested
		if err == nil && a.Normalize && res != nil {
			res.Normalize()
		}

		return res, err
	}
}
