package tema

// Segment represents the name of the segment
type Segment struct {
	ID string `json:"segment"`

	Hash    string `json:"hash"`    // the hash of this segment (if any)
	Touched bool   `json:"touched"` // has this segment been touched within recent changes
}

// SegmentMapping returns the mapping for the segments index
func (config *Configuration) SegmentMapping() interface{} {
	return j{
		"settings": j{
			"index": j{
				"number_of_replicas": 0,
				"number_of_shards":   128,
			},
		},
		"mappings": j{
			config.SegmentType: j{
				"dynamic": false,
				"properties": j{
					"segment": j{
						"type": "keyword",
					},
					"hash": j{
						"type": "keyword",
					},
					"touched": j{
						"type": "boolean",
					},
				},
			},
		},
	}
}
