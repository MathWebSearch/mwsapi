package result

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/MathWebSearch/mwsapi/utils"
	"github.com/pkg/errors"
)

// UnmarshalMWS un-marshals a resoponse from mws
func (res *Result) UnmarshalMWS(response *http.Response) error {
	defer response.Body.Close() // close the body when done
	decoder := json.NewDecoder(response.Body)

	// read the adapted data from the json
	var r struct {
		Total    int64 `json:"total"`
		TookInMS int   `json:"time"`

		Variables []*QueryVariable `json:"qvars"`

		MathWebSearchIDs []int64 `json:"ids,omitempty"`
		Hits             []*Hit  `json:"hits,omitempty"`
	}
	if err := decoder.Decode(&r); err != nil {
		err = errors.Wrap(err, "json.Unmarshal failed")
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
