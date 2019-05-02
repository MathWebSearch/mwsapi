package utils

import (
	"encoding/json"
	"fmt"
	"go/build"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"testing"
)

// CompareJSONAsset compares a result with a json file
func CompareJSONAsset(t *testing.T, name string, res interface{}, filename string) bool {
	// marshal the first file into json
	gbytes, err := json.Marshal(res)
	if err != nil {
		t.Errorf("%s Unable to marshal result: %s", name, err.Error())
		return false
	}
	got := string(gbytes)

	// load the file or fail
	ebytes, err := ioutil.ReadFile(filename)
	if err != nil {
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
		t.Errorf("%s json differs", name)
		return false
	}

	return true
}

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

// StartIntegrationTest starts the integration tests with the given name
func StartIntegrationTest(t *testing.M, name string) (code int) {
	return runBashScript(t, testDirPath, fmt.Sprintf("%s-up.sh", name))
}

// StopIntegrationTest stops the integration tests with the given name
func StopIntegrationTest(t *testing.M, name string) (code int) {
	return runBashScript(t, testDirPath, fmt.Sprintf("%s-down.sh", name))
}

// runBashScript runs a given script
func runBashScript(t *testing.M, pth string, script string) (code int) {

	// make a command
	cmd := exec.Command("/bin/bash", path.Join(pth, script))
	cmd.Dir = pth

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stdout

	cmd.Start()

	err := cmd.Wait()
	if err != nil {
		log.Print(err.Error())
		return 1
	}

	return 0
}

var testDirPath string

func init() {
	// get the base path
	p, err := build.Import("github.com/MathWebSearch/mwsapi", "", build.FindOnly)
	if err != nil {
		panic(err)
	}

	testDirPath = path.Join(p.Dir, "test")
}
