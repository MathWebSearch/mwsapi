package integrationtest

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/pkg/errors"
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
	filename = fmt.Sprintf("%s.out", asset)
	_, err = writeJSONFile(t, res, filename)
	err = errors.Wrap(err, "writeJSONFile failed")
	return
}

// writeJSONFile writes a json version of res into filename
func writeJSONFile(t *testing.T, res interface{}, filename string) (bytes []byte, err error) {
	// Remarshal it
	bytes, err = json.MarshalIndent(res, "", "  ")
	err = errors.Wrap(err, "json.MarshalIndent failed")
	if err != nil {
		return
	}

	// write the file
	err = ioutil.WriteFile(filename, bytes, 0755)
	err = errors.Wrap(err, "ioutil.WriteFile failed")
	return
}
