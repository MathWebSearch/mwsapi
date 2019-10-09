package utils

import (
	"bytes"

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
	if bytes.EqualFold(yesBytes, text) {
		*byesno = true
	} else if bytes.EqualFold(noBytes, text) {
		*byesno = false

		// do not load the else
	} else {
		err = errors.Errorf("Boolean should be \"yes\" or \"no\", not %q", string(text))
	}

	return
}

var noBytes []byte
var yesBytes []byte

func init() {
	noBytes = []byte("no")
	yesBytes = []byte("yes")
}
