package elasticsync

import (
	"github.com/MathWebSearch/mwsapi/connection"
)

type j map[string]interface{}

// harvestMapping returns the mapping for the harvest index
func harvestMapping(config *connection.ElasticConfiguration) interface{} {
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

// segmentMapping returns the mapping for the segments index
func segmentMapping(config *connection.ElasticConfiguration) interface{} {
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
