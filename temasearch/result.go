package temasearch

import (
	"time"

	"github.com/MathWebSearch/mwsapi/mws"
	"github.com/MathWebSearch/mwsapi/tema/query"
)

// Result represents a single MathWebSearch result
type Result struct {
	Total int64 `json:"total"` // the total number of results
	From  int64 `json:"from"`  // result number this page starts at
	Size  int64 `json:"size"`  // (maximum) number of results in this page

	Took *time.Duration `json:"took"` // the amount of time the query took to execute, including network latency to elasticsearch

	Variables []*QueryVariable `json:"qvars"` // list of query variables in the original query
	Hits      []*Hit           `json:"hits"`  // the current page of hits
}

// Hit represents a single MathWebSearch Hit
type Hit struct {
	ID       string      `json:"id"`       // the id of the document returned
	Metadata interface{} `json:"metadata"` // the metadata of the hit
	XHTML    string      `json:"xhtml"`    // returned xhtml (if any)

	Score    float64  `json:"score"`    // the score this hit achieved
	Snippets []string `json:"snippets"` // the generted snippets

	Math []*ReplacedFormulaInfo `json:"maths"` // the replaced math elements within this snippet

}

// QueryVariable represents a query variable
type QueryVariable struct {
	*mws.QueryVariable
}

// ReplacedFormulaInfo represents a single replaced Except
type ReplacedFormulaInfo struct {
	*query.ReplacedFormulaInfo
}

func (res *Result) fromMathWebSearch(mws *mws.Result) {
	// transform the results
	res.Size = mws.Size
	res.Total = mws.Total

	// copy over variables
	res.Variables = make([]*QueryVariable, len(mws.Variables))
	for i, e := range mws.Variables {
		res.Variables[i] = &QueryVariable{e}
	}

	// copy over the hits
	res.Hits = make([]*Hit, len(mws.Hits))
	for i, e := range mws.Hits {
		res.Hits[i] = newHitFromMWS(e)
	}

}
func (res *Result) fromElastic(elastic *query.Result) {
	// transform the results
	res.Total = elastic.Total
	res.Size = elastic.Size

	// copy over the hits
	res.Hits = make([]*Hit, len(elastic.Hits))
	for i, e := range elastic.Hits {
		res.Hits[i] = newHitFromElastic(e)
	}

	return
}

func newHitFromMWS(mws *mws.Hit) (hit *Hit) {
	// copy over the xhtml
	hit = &Hit{
		XHTML: mws.XHTML,
	}

	// and the replaced formulae info
	hit.Math = make([]*ReplacedFormulaInfo, len(mws.Formulae))
	for i, e := range mws.Formulae {
		hit.Math[i] = &ReplacedFormulaInfo{
			&query.ReplacedFormulaInfo{
				URL:   e.URL,
				XPath: e.XPath,
			},
		}
	}

	return
}

func newHitFromElastic(tema *query.Hit) (hit *Hit) {
	hit = &Hit{
		ID:       tema.ElasticID,
		Metadata: tema.Metadata,

		Score:    tema.Score,
		Snippets: tema.Snippets,
	}

	// and the replaced formulae info
	hit.Math = make([]*ReplacedFormulaInfo, len(tema.Math))
	for i, e := range tema.Math {
		hit.Math[i] = &ReplacedFormulaInfo{e}
	}

	return
}
