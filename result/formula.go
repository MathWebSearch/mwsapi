package result

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/MathWebSearch/mwsapi/utils"
)

// MathFormula represents information about a single located math formula
// It supports JSON Marshal + Unmarshal (for reading it from mws results and sending it to the user)
// And XML Un-Marshelling (for reading it as part of raw harvest data)
type MathFormula struct {
	Source string `json:"source,omitempty" xml:",chardata"` // MathML Element Representing entire formula

	DocumentURL string `json:"durl,omitempty"`          // document url (if any)
	LocalID     string `json:"url" xml:"local_id,attr"` // local formula id

	XPath   string `json:"xpath,omitempty"`   // XPath from the formula -> query
	SubTerm string `json:"subterm,omitempty"` // MathML Element representing matching subterm

	Substitution map[string]string `json:"subst,omitempty"` // MathML Elements representing values for the subsituted terms
}

// MathML returns a MathML object representing this MathFormula
func (formula *MathFormula) MathML() (*utils.MathML, error) {
	if formula.Source == "" {
		return nil, errors.New("[MathFormula.MathML] No Source available")
	}
	return utils.ParseMathML(formula.Source)
}

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

// RealMathID returns the math id used for this object in dictionaries
func (formula *MathFormula) RealMathID() string {
	mathid := formula.LocalID
	if _, err := strconv.Atoi(mathid); err == nil {
		return "math" + mathid
	}
	return mathid
}

// PopulateSubsitutions populates the subsitutions field of this MathFormula
// given a hit and a result
func (formula *MathFormula) PopulateSubsitutions(hit *Hit, res *Result) (err error) {
	if len(formula.Substitution) > 0 {
		return fmt.Errorf("[MathFormula.PopulateSubsitutions] Substiution already populated ")
	}

	if formula.Source == "" {
		return fmt.Errorf("[MathFormula.PopulateSubsitutions] Missing formula source, can not populate subsitutions")
	}

	// parse the mathml
	mathml, err := formula.MathML()
	if err != nil {
		return
	}

	// find the term representing the entire found term
	err = mathml.NavigateAnnotation(".."+formula.XPath, false)
	if err != nil {
		return err
	}

	// store the subterm that matches
	formula.SubTerm = mathml.OutputXML()

	formula.Substitution = make(map[string]string, len(res.Variables))

	// iterate over the variables
	for _, variable := range res.Variables {
		// make a copy of the mathml object
		copy := mathml.Copy()

		// navigate to the xpath
		err = copy.NavigateAnnotation("."+variable.XPath, true)
		if err != nil {
			return err
		}

		// and output the xml
		formula.Substitution[variable.Name] = copy.OutputXML()
	}

	return
}
