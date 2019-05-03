package integrationtest

import (
	"os/exec"
	"testing"
)

// MarkSkippable marks t as containing a skippable integration test
// returns true if the test has been skipped
func MarkSkippable(t *testing.T) bool {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
		return true
	}
	if !hasDockerAndCompose() {
		t.Skip("skipping integration test due to incorrect docker-compose setup")
		return true
	}

	return false
}

// hasDockerAndCompose checks if the 'docker' and 'docker-compose' executables are installed
func hasDockerAndCompose() bool {
	if _, err := exec.LookPath("docker"); err != nil {
		return false
	}

	if _, err := exec.LookPath("docker-compose"); err != nil {
		return false
	}

	return true
}
