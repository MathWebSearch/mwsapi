package result

import (
	"strconv"
	"strings"
)

// MathInfo represents information about a math element within a result
type MathInfo struct {
	Source string `json:"source,omitempty"` // html source code of this formula (if any)
	URL    string `json:"url"`              // URL of the replaced formula
	XPath  string `json:"xpath"`            // xpath of the term within the formula
}

// ID returns the (local) id for this MathInfo Element
func (info *MathInfo) ID() string {
	parts := strings.Split(info.URL, "#")
	return parts[len(parts)-1]
}

// RealMathID return the real math id
func (info *MathInfo) RealMathID() string {
	mathid := info.ID()
	if _, err := strconv.Atoi(mathid); err == nil {
		return "math" + mathid
	}
	return mathid
}

// Variable represents a query variable
type Variable struct {
	Name string `json:"name"` // the name of this query variable
}
