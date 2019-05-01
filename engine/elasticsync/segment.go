package elasticsync

// Segment represents the name of the segment
type Segment struct {
	ID string `json:"segment"`

	Hash    string `json:"hash"`    // the hash of this segment (if any)
	Touched bool   `json:"touched"` // has this segment been touched within recent changes
}
