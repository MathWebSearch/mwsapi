package mwsengine

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/MathWebSearch/mwsapi/result"
	"github.com/MathWebSearch/mwsapi/utils"
)

// rawResult represents a raw MWS Result
type rawResult struct {
	Total int64 `json:"total"` // total size of the resultset (if requested)

	TookInMS int `json:"time"` // how long the query took, in ms

	Variables []*result.Variable `json:"qvars"` // list of query variables in the original query

	MathWebSearchIDs []int64       `json:"ids,omitempty"`  // MathWebSearchIds
	Hits             []*result.Hit `json:"hits,omitempty"` // the actual hits of this element
}

// newMWSResult populates a non-nil result with mws results
func newMWSResult(res *result.Result, response *http.Response) (err error) {
	defer response.Body.Close()

	// read the raw result
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return
	}
	var raw rawResult
	err = json.Unmarshal(body, &raw)

	took := time.Duration(raw.TookInMS) * time.Millisecond
	res.Took = &took

	res.Kind = "mwsd"
	res.Total = raw.Total

	res.TookComponents = map[string]*time.Duration{
		"mwsd": &took,
	}

	res.Size = utils.MaxInt64(int64(len(raw.Hits)), int64(len(raw.MathWebSearchIDs)))

	res.Variables = raw.Variables
	res.HitIDs = raw.MathWebSearchIDs
	res.Hits = raw.Hits

	return
}
