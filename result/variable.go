package result

// QueryVariable represents a query variable
type QueryVariable struct {
	Name  string `json:"name"`  // name of the variable
	XPath string `json:"xpath"` // xpath of the variable relative to the root
}
