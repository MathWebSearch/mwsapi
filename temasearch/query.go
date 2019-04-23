package temasearch

// Query is a TemaSearch Query
type Query struct {
	Expressions []string // a list of expressions to search
	Text        string   // some MathWebSearch Text to search for
}

// HasText checks if a query has text
func (q *Query) HasText() bool {
	return q.Text != ""
}

// HasExpressions checks if a query has expressions
func (q *Query) HasExpressions() bool {
	return len(q.Expressions) > 0
}
