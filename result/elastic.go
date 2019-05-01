package result

// ElasticElement represents a single element within an elasticsearch harvest
type ElasticElement struct {
	Metadata interface{} `json:"metadata"`          // Arbitrary document metadata
	Segment  string      `json:"segment,omitempty"` // Name of the segment this document belongs to

	Text string `json:"text"` // Text of this element

	MWSPaths   map[int64]map[string]MathInfo `json:"mws_id"`  // information about each identifier within this document
	MWSNumbers []int64                       `json:"mws_ids"` // list of identifiers within this document

	MathSource map[string]string `json:"math"` // Source of replaced math elements within this document
}

// ElasticSegment represents the name of the segment
type ElasticSegment struct {
	ID string `json:"segment"`

	Hash    string `json:"hash"`    // the hash of this segment (if any)
	Touched bool   `json:"touched"` // has this segment been touched within recent changes
}
