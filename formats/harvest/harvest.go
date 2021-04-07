// Package harvest implements the MathWebSearch Harvest format.
package harvest

// Harvest represents a single MathWebSearch Harvest.
type Harvest struct {
	Data        []HarvestData
	Expressions []HarvestExpr
}

// HarvestData represents information about a Data Element of MathWebSearch.
type HarvestData struct {
	ID      string // ID is the ID of this Harvest:Data
	Content string // Content is the Content of this Harvest:Data Element
}

// HarvestExpr represents information about a Harvest Expression
type HarvestExpr struct {
	ID     string // ID of this expression
	DataID string // the ID of the data element belonging to this expression
	Math   string // The corresponding math element to this Harvest expression
}
