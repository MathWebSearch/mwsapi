package integrationtest

import (
	"fmt"
	"io/ioutil"
	"testing"

	jsoniter "github.com/json-iterator/go"

	"github.com/pkg/errors"
)

// normalizeJSON normalizes json into an interface{}
func normalizeJSON(in string) (interface{}, error) {
	var t interface{}

	// Unmarshal into a generic interface
	err := jsoniter.Unmarshal([]byte(in), &t)
	if err != nil {
		return nil, err
	}

	return t, nil
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
	bytes, err = jsoniter.MarshalIndent(res, "", "  ")
	err = errors.Wrap(err, "jsoniter.MarshalIndent failed")
	if err != nil {
		return
	}

	// write the file
	err = ioutil.WriteFile(filename, bytes, 0755)
	err = errors.Wrap(err, "ioutil.WriteFile failed")
	return
}
