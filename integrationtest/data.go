package integrationtest

import (
	"go/build"
	"os"
	"path"
)

// TestDataPath contains the full path to the 'testdata' directory
// This variable is guaranteed to be set at runtime; if the directory does not exist panic() is called
var TestDataPath string

func init() {
	// build the path to the testdata directory
	p, _ := build.Import("github.com/MathWebSearch/mwsapi/integrationtest", "", build.FindOnly)
	TestDataPath = path.Join(p.Dir, "testdata")

	// panic if it does not exist
	_, err := os.Stat(TestDataPath)
	if err != nil && os.IsNotExist(err) {
		panic("Integrationtest testdata not found")
	}
}
