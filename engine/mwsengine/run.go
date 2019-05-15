package mwsengine

import (
	"bytes"
	"encoding/xml"
	"net/http"
	"time"

	"github.com/MathWebSearch/mwsapi/connection"
	"github.com/MathWebSearch/mwsapi/query"
	"github.com/MathWebSearch/mwsapi/result"
)

// Run runs an MWS Query
func Run(conn *connection.MWSConnection, query *query.MWSQuery, from int64, size int64) (res *result.Result, err error) {
	// measure time for this query
	start := time.Now()
	defer func() {
		if res != nil {
			took := time.Since(start)
			res.Took = &took
		}
	}()

	// TODO: Paralellize this with appropriate page size
	return runRaw(conn, query.Raw(from, size))
}

// RunRaw runs a raw query
func runRaw(conn *connection.MWSConnection, q *query.RawMWSQuery) (res *result.Result, err error) {
	// TODO: Split this into smaller queries of at most size PageSize
	// and then join all of them together

	// turn the query into xml
	b, err := xml.Marshal(q)
	if err != nil {
		return
	}

	// make a request object
	req, err := http.NewRequest("POST", conn.URL(), bytes.NewBuffer(b))
	if err != nil {
		return
	}

	// set some headers
	req.Header.Set("Content-Type", "application/xml")

	// run the request
	resp, err := conn.Client.Do(req)
	if err != nil {
		return
	}

	// initialize the result
	res = &result.Result{
		From: q.From,
		Size: q.Size,
	}
	err = res.UnmarshalMWS(resp)

	return
}
