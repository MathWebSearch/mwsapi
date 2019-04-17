package tema

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"os"
)

// Segment represents the name of the segment
type Segment struct {
	ID string `json:"segment"`

	Hash    string `json:"hash"`    // the hash of this segment (if any)
	Touched bool   `json:"touched"` // has this segment been touched within recent changes
}

// HashSegment computes the hash of a segment
func HashSegment(filename string) (hash string, err error) {
	// the hasher implementation
	hasher := sha256.New()

	// open the segment
	f, err := os.Open(filename)
	if err != nil {
		return
	}
	defer f.Close()

	// start hashing the file
	if _, err := io.Copy(hasher, f); err != nil {
		return "", err
	}

	// turn it into a string
	hash = hex.EncodeToString(hasher.Sum(nil))
	return
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
