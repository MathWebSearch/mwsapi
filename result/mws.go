package result

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/MathWebSearch/mwsapi/utils"
)

// UnmarshalMWS un-marshals a resoponse from mws
func (res *Result) UnmarshalMWS(response *http.Response) error {
	defer response.Body.Close() // close the body when done
	responseBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}

	// read the adapted data from the json
	var r struct {
		Total    int64 `json:"total"`
		TookInMS int   `json:"time"`

		Variables []*QueryVariable `json:"qvars"`

		MathWebSearchIDs []int64 `json:"ids,omitempty"`
		Hits             []*Hit  `json:"hits,omitempty"`
	}
	if err := json.Unmarshal(responseBytes, &r); err != nil {
		return err
	}

	// store the time taken
	took := time.Duration(r.TookInMS) * time.Millisecond
	res.Took = &took

	res.Kind = MathWebSearchKind
	res.Total = r.Total

	// TODO: Update this
	res.TookComponents = map[string]*time.Duration{
		"mwsd": &took,
	}

	res.Size = utils.MaxInt64(int64(len(r.Hits)), int64(len(r.MathWebSearchIDs)))

	res.Variables = r.Variables
	res.HitIDs = r.MathWebSearchIDs
	res.Hits = r.Hits

	return res.PopulateSubsitutions()
}
