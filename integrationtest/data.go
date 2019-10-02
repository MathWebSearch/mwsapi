package integrationtest

import (
	"os"
	"path"
	"runtime"
)

// TestDataPath contains the full path to the 'testdata' directory
// This variable is guaranteed to be set at runtime; if the directory does not exist panic() is called
var TestDataPath string

func init() {
	// build the path to the testdata directory
	_, p, _, ok := runtime.Caller(0)
	if !ok {
		panic("runtime.caller failed")
	}
	TestDataPath = path.Join(p, "..", "testdata")

	// panic if it does not exist
	_, err := os.Stat(TestDataPath)
	if err != nil && os.IsNotExist(err) {
		panic("Integrationtest testdata not found")
	}
}
