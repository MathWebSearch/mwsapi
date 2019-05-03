package integrationtest

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"testing"
)

// normalizeJSON normalizes a json string
func normalizeJSON(in string) (out string, err error) {
	var t interface{}

	// Unmarshal into a generic interface
	err = json.Unmarshal([]byte(in), &t)
	if err != nil {
		return
	}

	// Remarshal it
	outB, err := json.Marshal(&t)
	if err != nil {
		return
	}

	// and return
	out = string(outB)
	return
}

// Outputs json into {$asset}.out for debugging purposes
// returns the filename
func outputDebugJSON(t *testing.T, res interface{}, asset string) (filename string, err error) {
	// filename to place broken output into
	filename = fmt.Sprintf("%s.out", asset)

	// Remarshal it
	bytes, err := json.MarshalIndent(res, "", "  ")
	if err != nil {
		return
	}

	// write the file
	err = ioutil.WriteFile(filename, bytes, 0755)
	if err != nil {
		return
	}

	return
}
