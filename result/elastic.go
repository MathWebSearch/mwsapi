package result

import (
	"time"

	"github.com/MathWebSearch/mwsapi/utils/elasticutils"
	"github.com/pkg/errors"
)

// UnmarshalElastic un-marshals a response from elastic
func (res *Result) UnmarshalElastic(page *elasticutils.ResultsPage) error {
	// it is an elasticsearch document query
	res.Kind = ElasticDocumentKind

	// prepare result objejct
	res.Total = page.Total
	res.Hits = make([]*Hit, len(page.Hits))
	res.Size = int64(len(page.Hits))

	// and make the new hits
	for i, hit := range page.Hits {
		res.Hits[i] = &Hit{}
		if err := res.Hits[i].UnmarshalElasticDocument(hit); err != nil {
			return errors.Wrap(err, "res.Hits[i].UnmarshalElasticDocument failed")
		}
		if err := res.Hits[i].PopulateSubsitutions(res); err != nil {
			return errors.Wrap(err, "res.Hits[i].PopulateSubsitutions failed")
		}

	}

	// add the time it took in the document phase
	document := time.Duration(page.Took) * time.Millisecond
	res.TookComponents = map[string]*time.Duration{
		"document": &document,
	}

	return nil
}

// UnmarshalElasticDocument unmarshals a document hit from elasticsearch
func (hit *Hit) UnmarshalElasticDocument(obj *elasticutils.Object) (err error) {
	// create the Hit and set it's id properly
	hit.ID = obj.GetID()

	// un-marshal the harvest element
	err = obj.Unpack(&hit.Element)
	err = errors.Wrap(err, "obj.Unpack failed")
	if err != nil {
		return
	}

	// and populate the math field
	return hit.PopulateMath()
}

// UnmarshalElasticHighlight populates a document hit with a highlight hit
func (hit *Hit) UnmarshalElasticHighlight(obj *elasticutils.Object) (err error) {
	if obj.Hit == nil || obj.Hit.Highlight == nil {
		return errors.New("[Hit.UnmarshalElasticHighlight] No highlights returned")
	}

	// load the highlights
	var ok bool
	hit.Snippets, ok = (*obj.Hit.Highlight)["text"]
	if !ok {
		return errors.New("[Hit.UnmarshalElasticHighlight] No highlights returned")
	}

	return
}
