package tema

// HarvestElement represents a single element of a harvest
type HarvestElement struct {
	Metadata interface{} `json:"metadata"`          // Arbitrary document metadata
	Segment  string      `json:"segment,omitempty"` // Name of the segment this document belongs to

	MWSPaths   map[int64]map[string]PathInfo `json:"mws_id"`  // information about each identifier within this document
	MWSNumbers []int64                       `json:"mws_ids"` // list of identifiers within this document

	MathSource map[string]string `json:"math"` // Source of replaced math elements within this document
}

// PathInfo represents information about a subterm
type PathInfo struct {
	XPath string `json:"xpath"`
}

// HarvestMapping returns the mapping for the harvest index
func (config *Configuration) HarvestMapping() interface{} {
	return j{
		"settings": j{
			"index": j{
				"number_of_replicas": 0,
				"number_of_shards":   128,
			},
		},
		"mappings": j{
			config.HarvestType: j{
				"dynamic": false,
				"properties": j{
					"metadata": j{
						"dynamic": true,
						"type":    "object",
					},
					"segment": j{
						"type": "keyword",
					},
					"mws_ids": j{
						"type":  "long",
						"store": false,
					},
					"text": j{
						"type":  "text",
						"store": false,
					},
					"mws_id": j{
						"enabled": false,
						"type":    "object",
					},
					"math": j{
						"enabled": false,
						"type":    "object",
					},
				},
			},
		},
	}
}
