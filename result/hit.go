package result

// Hit represents a single Hit
type Hit struct {
	ID string `json:"id,omitempty"` // the (possibly internal) id of this hit

	URL   string `json:"url,omitempty"` // the url of the document returned
	XPath string `json:"xpath"`         // the xpath of the query term to the formulae referred to by this id

	Element *HarvestElement `json:"source,omitempty"` // the raw ElasticSearch element (if any)

	Metadata interface{} `json:"metadata,omitempty"` // arbitrary document meta-data

	Score float64 `json:"score,omitempty"` // score of this hit

	Snippets []string `json:"snippets,omitempty"` // extracts of this hit (if any)
	XHTML    string   `json:"xhtml,omitempty"`    // xhtml source of this hit (if available)

	Math []*MathFormula `json:"math_ids"` // math found within this hit
}
