package integrationtest

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"testing"
)

// TestJSONAsset compares a result with a json file
func TestJSONAsset(t *testing.T, name string, res interface{}, filename string) bool {
	// marshal the first file into json
	gbytes, err := json.Marshal(res)
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
	// Read json, then re-marshal

	// trim all the things
	got, err = normalizeJSON(got)
	if err != nil {
		t.Errorf("%s got invalid json: %s", name, err.Error())
		return false
	}

	expected, err = normalizeJSON(expected)
	if err != nil {
		t.Errorf("%s invalid json in asset %q: %s", name, filename, err.Error())
		return false
	}

	if got != expected {
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
