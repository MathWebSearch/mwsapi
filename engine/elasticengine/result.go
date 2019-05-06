package elasticengine

import (
	"errors"
	"fmt"
	"sort"
	"time"

	"github.com/MathWebSearch/mwsapi/result"
	"github.com/MathWebSearch/mwsapi/utils/elasticutils"
)

// newDocumentResult populates a non-nil result by using results form a page
func newDocumentResult(res *result.Result, page *elasticutils.ResultsPage) (err error) {

	// it is an elasticsearch document query
	res.Kind = "elastic-document"

	// prepare result objejct
	res.Total = page.Total
	res.Hits = make([]*result.Hit, len(page.Hits))
	res.Size = int64(len(page.Hits))

	// and make the new hits
	for i, hit := range page.Hits {
		res.Hits[i], err = newDocumentHit(hit)
		if err != nil {
			return err
		}
	}

	res.TookComponents = map[string]*time.Duration{
		"document": &page.Took,
	}

	return
}

// newDocumentHit creates a new Hit within a document result
func newDocumentHit(obj *elasticutils.Object) (hit *result.Hit, err error) {
	hit = &result.Hit{
		ID: obj.GetID(),
	}

	// unpack the result element we get from elastic
	var raw result.ElasticElement
	err = obj.Unpack(&raw)
	if err != nil {
		return
	}
	hit.Element = &raw

	// create the math elements, without knowing the size beforehand
	hit.Math = []*result.MathInfo{}

	for _, mwsid := range raw.MWSNumbers {
		// load the data
		data, ok := raw.MWSPaths[mwsid]
		if !ok {
			return nil, fmt.Errorf("Result %q missing path info for %d", hit.ID, mwsid)
		}

		// sort the keys in alphabetical order
		keys := make([]string, len(data))
		for key := range data {
			keys = append(keys, key)
		}
		sort.Strings(keys)

		// and iterate over it
		for _, key := range keys {
			value := data[key]
			hit.Math = append(hit.Math, &result.MathInfo{
				URL:   key,
				XPath: value.XPath,
			})
		}
	}

	return
}

// newHighlightHit populates a document hit with a highlight hit
func newHighlightHit(hit *result.Hit, obj *elasticutils.Object) (err error) {
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
		hit.Math[i].Source, ok = hit.Element.MathSource[math.ID()]
		if !ok {
			return fmt.Errorf("Result %s with source info %#v missing info for %s", hit.ID, hit.Element.MathSource, math.ID())
		}
	}

	return
}
