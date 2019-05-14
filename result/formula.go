package result

import (
	"strconv"
	"strings"
)

// MathFormula represents information about a single located math formula
// It supports JSON Marshal + Unmarshal (for reading it from mws results and sending it to the user)
// And XML Un-Marshelling (for reading it as part of raw harvest data)
type MathFormula struct {
	Source string `json:"source,omitempty" xml:"innerxml"` // html source code of this formula (if any)

	DocumentURL string `json:"durl,omitempty" xml:"-"`  // document url (if any)
	LocalID     string `json:"url" xml:"local_id,attr"` // local formula id

	XPath        string            `json:"xpath,omitempty" xml:"-"` // XPath from the formula -> query
	Substitution map[string]string `json:"subst,omitempty" xml:"-"` // values for the replaced terms
}

// TODO: Rename local-id field in json
// TODO: Add JSON + XML Tests
// TODO: Check if we need the documentURL

// SetURL sets the url of a MathInfo object
func (formula *MathFormula) SetURL(url string) {
	idx := strings.LastIndex(url, "#")

	// if we have no '#' we only have a local id
	if idx == -1 {
		formula.DocumentURL = ""
		formula.LocalID = url
		return
	}

	// set the appropriate parts of the url
	formula.DocumentURL = url[:idx]
	formula.LocalID = url[idx+1:]
}

// RealMathID return the real math id
func (formula *MathFormula) RealMathID() string {
	mathid := formula.LocalID
	if _, err := strconv.Atoi(mathid); err == nil {
		return "math" + mathid
	}
	return mathid
}
