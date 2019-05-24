package utils

import (
	"strings"

	"github.com/pkg/errors"
)

// BooleanYesNo represents a boolean that is xml encoded as "yes" or "no"
type BooleanYesNo bool

// MarshalText turns a BooleanYesNo into a string
func (byesno BooleanYesNo) MarshalText() (text []byte, err error) {
	if byesno {
		text = []byte("yes")
	} else {
		text = []byte("no")
	}
	return
}

// UnmarshalText unmarshals text into a string
func (byesno *BooleanYesNo) UnmarshalText(text []byte) (err error) {
	// load yes and no
	qtext := strings.ToLower(string(text))
	if qtext == "yes" {
		*byesno = true
	} else if qtext == "no" {
		*byesno = false

		// do not load the else
	} else {
		err = errors.Errorf("Boolean should be \"yes\" or \"no\", not %q", qtext)
	}

	return
}
