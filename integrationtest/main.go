package integrationtest

import (
	"flag"
	"os"
	"testing"
)

// Main is the entry point for an integration test
// service and urls are passed directly to StartDockerService and StopDockerService
func Main(m *testing.M, service string, urls ...string) {
	// parse the flags manually to ensure that .Verbose() and .Short() are set
	flag.Parse()

	var code int
	defer os.Exit(code)

	// start the docker service
	if err := StartDockerService(service, urls...); err != nil {
		panic(err)
	}
	defer StopDockerService(service)

	code = m.Run()
}
