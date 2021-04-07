package harvest

import (
	"errors"
	"fmt"
)

// Validate checks if this harvest is valid.
//
// When a harvest is not valid, it returns an error of one of the ErrValidate* types (or constants) in this package.
// When a harvest is valid, it returns nil.
func (h Harvest) Validate() error {

	// check that there is at least one expression!
	if len(h.Expressions) == 0 {
		return ErrValidateNoExpressions
	}

	// check that each Data ID occurs exactly once
	// store all valid Data IDs in a map for easy checks later
	dataIDs := make(map[string]struct{}, len(h.Data))
	for _, data := range h.Data {
		if _, ok := dataIDs[data.ID]; ok {
			return ErrValidateDuplicateDataID{ID: data.ID}
		}
		dataIDs[data.ID] = struct{}{}
	}

	// check that each expression ID occurs exactly once
	// check that each referenced Data ID is valid (by checking the map above)
	exprIDs := make(map[string]struct{}, len(h.Expressions))
	for _, expr := range h.Expressions {
		if _, ok := exprIDs[expr.ID]; ok {
			return ErrValidateDuplicateExprID{ID: expr.ID}
		}
		if _, ok := dataIDs[expr.DataID]; !ok {
			return ErrValidateMissingData{ExprID: expr.ID, DataID: expr.DataID}
		}
		exprIDs[expr.ID] = struct{}{}
	}

	// TODO: validate that all the math is actually a valid math element!

	return nil
}

// ErrValidateNoExpressions indicates that the harvest contains no expressions
var ErrValidateNoExpressions = errors.New("Harvest contains no expressions")

// ErrValidateDuplicateDataID indicates that a data ID occurs twice
type ErrValidateDuplicateDataID struct {
	ID string
}

func (e ErrValidateDuplicateDataID) Error() string { return fmt.Sprintf("Duplicate Data ID: %q", e.ID) }

// ErrDuplicateDataID indicates that a data ID occurs twice
type ErrValidateDuplicateExprID struct {
	ID string
}

func (e ErrValidateDuplicateExprID) Error() string {
	return fmt.Sprintf("Duplicate Expression ID: %q", e.ID)
}

// ErrValidateMissingData indicates that an expression is referencing a missing Data ID
type ErrValidateMissingData struct {
	ExprID string
	DataID string
}

func (e ErrValidateMissingData) Error() string {
	return fmt.Sprintf("Expression %q is referencing missing Data ID %q", e.ExprID, e.DataID)
}
