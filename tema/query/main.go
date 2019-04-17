package query

// Query represents a query sent to temasearch
type Query struct {
	// list of possible formula ids
	FormulaID []string

	// Textual content
	Text string
}
