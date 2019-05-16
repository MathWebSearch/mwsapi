package result

import (
	"errors"

	"github.com/MathWebSearch/mwsapi/utils"
)

// QueryVariable represents a query variable
type QueryVariable struct {
	Name  string `json:"name"`  // name of the variable
	XPath string `json:"xpath"` // xpath of the variable relative to the root
}

// Value finds the value of this queryvariable within a given term
func (qvar *QueryVariable) Value(resultTerm string) (value string, err error) {
	// find all the matches
	values, err := utils.ResolveXPath(resultTerm, qvar.XPath)
	if err != nil {
		return
	}

	// check that we have a result
	if len(values) == 0 {
		err = errors.New("[QueryVariable.Value] XPath did not find any results")
		return
	}

	// and return it
	value = values[0]
	return
}
