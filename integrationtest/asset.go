package integrationtest

import (
	"io/ioutil"
	"os"
	"reflect"
	"testing"

	jsoniter "github.com/json-iterator/go"
)

// TestJSONAsset compares a result with a json file
func TestJSONAsset(t *testing.T, name string, res interface{}, filename string) bool {
	// marshal the first file into json
	gbytes, err := jsoniter.Marshal(res)
	if err != nil {
		t.Errorf("%s Unable to marshal result: %s", name, err.Error())
		return false
	}
	got := string(gbytes)

	// load the file or fail
	ebytes, err := ioutil.ReadFile(filename)

	// file does not exist => create an empty asset
	if err != nil && os.IsNotExist(err) {
		t.Errorf("%s Unable to load asset %q, creating file with 'null' in it. ", name, filename)
		ebytes, err = writeJSONFile(t, nil, filename)
		if err != nil {
			t.Errorf("%s Unable to write asset %q: %s", name, filename, err.Error())
			return false
		}

		// else just throw the error message
	} else if err != nil {
		t.Errorf("%s Unable to load asset %q: %s", name, filename, err.Error())
		return false
	}
	expected := string(ebytes)

	gotJSON, err := normalizeJSON(got)
	if err != nil {
		return false
	}
	expectedJSON, err := normalizeJSON(expected)
	if err != nil {
		return false
	}

	// Read json, then re-marshal
	if !reflect.DeepEqual(gotJSON, expectedJSON) {
		fn, err := outputDebugJSON(t, res, filename)

		if err == nil {
			t.Errorf("%s json differs, wrote output in %q", name, fn)
		} else {
			t.Errorf("%s json differs, but an error occured while trying to write output in %q: %s", name, fn, err.Error())
		}
		return false
	}

	return true
}
