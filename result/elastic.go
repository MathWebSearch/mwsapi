package result

import (
	"errors"
	"fmt"
	"sort"
	"time"

	"github.com/MathWebSearch/mwsapi/utils/elasticutils"
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
		err := res.Hits[i].UnmarshalElasticDocument(hit)
		if err != nil {
			return err
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
	*hit = Hit{
		ID: obj.GetID(),
	}

	// un-marshal the harvest element
	err = obj.Unpack(&hit.Element)
	if err != nil {
		return
	}

	// load all the MWSElement and the paths
	// TODO: Move this into a seperate method later on
	for _, mwsid := range hit.Element.MWSNumbers {
		// load the data
		data, ok := hit.Element.MWSPaths[mwsid]
		if !ok {
			return fmt.Errorf("Result %q missing path info for %d", hit.ID, mwsid)
		}

		// sort the keys in alphabetical order
		keys := make([]string, len(data))
		count := 0
		for key := range data {
			keys[count] = key
			count++
		}
		sort.Strings(keys)

		// and iterate over it
		for _, key := range keys {
			res := &MathFormula{XPath: data[key].XPath}
			res.SetURL(key)
			hit.Math = append(hit.Math, res)
		}
	}

	return
}

// UnmarshalElasticHighlight populates a document hit with a highlight hit
func (hit *Hit) UnmarshalElasticHighlight(obj *elasticutils.Object) (err error) {
	if obj.Hit == nil || obj.Hit.Highlight == nil {
		return errors.New("No highlights returned")
	}

	// load the highlights
	var ok bool
	hit.Snippets, ok = (*obj.Hit.Highlight)["text"]
	if !ok {
		return errors.New("No highlights returned")
	}

	// map() over doc.Math
	for i, math := range hit.Math {
		var ok bool
		hit.Math[i].Source, ok = hit.Element.MathSource[math.LocalID]
		if !ok {
			return fmt.Errorf("Result %s with source info %#v missing info for %s", hit.ID, hit.Element.MathSource, math.LocalID)
		}
	}

	return
}
