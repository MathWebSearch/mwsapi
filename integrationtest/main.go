package integrationtest

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"testing"
)

// Main is the entry point for an integration test
// service and urls are passed directly to StartDockerService and StopDockerService
func Main(m *testing.M, service string, init func(client *http.Client) error, urls ...string) {
	// parse the flags manually to ensure that .Verbose() and .Short() are set
	flag.Parse()

	var code int
	defer os.Exit(code)

	// start the docker service
	client, err := StartDockerService(service, urls...)
	if err != nil {
		panic(err)
	}
	defer StopDockerService(service)

	// run init code
	if init != nil {
		if err := init(client); err != nil {
			fmt.Println(err.Error())
			code = 1
			return
		}
	}

	code = m.Run()
}
