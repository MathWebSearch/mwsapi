package utils

import (
	"fmt"
	"go/build"
	"log"
	"os"
	"os/exec"
	"path"
	"testing"
)

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
