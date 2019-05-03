package integrationtest

import (
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"sync"
	"testing"
	"time"
)

// StartDockerService starts a docker-compose configured service if the short tests are not set
// and then blocks until a url returns http status code 200
func StartDockerService(service string, urls ...string) (err error) {
	if testing.Short() {
		return
	}

	err = runExternalCommand("docker-compose", "up", "--force-recreate", "-d", service)
	if err != nil {
		return
	}

	// create a new client and wait group
	client := &http.Client{}

	group := &sync.WaitGroup{}
	group.Add(len(urls))

	for _, url := range urls {
		go func(url string) {
			defer group.Done()
			for {
				if testing.Verbose() {
					fmt.Printf("Checking if %q returns HTTP 200 ...\n", url)
				}
				res, err := client.Get(url)
				if err == nil && res.StatusCode == 200 {
					break
				}
				time.Sleep(1 * time.Second) // wait for the next one
			}
		}(url)
	}

	// wait for all of them to be done
	group.Wait()
	return
}

// StopDockerService gracefully stops and then removes a docker-compose configured service if the short flag is not set
// and all of it's associated volumes
func StopDockerService(service string) (err error) {
	if testing.Short() {
		return
	}

	// stop the service
	err = runExternalCommand("docker-compose", "stop", service)
	if err != nil {
		return
	}

	// remove it and its (anonymous volumes)
	return runExternalCommand("docker-compose", "rm", "-f", "-v", service)
}

// runExternalCommand runs an external command in the TestDataPath directory
func runExternalCommand(exe string, args ...string) (err error) {
	// set command, path + arguments
	cmd := exec.Command(exe, args...)
	cmd.Dir = TestDataPath

	// if in verbose mode
	// show a bunch of output to the user
	if testing.Verbose() {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}

	// finally run the command
	return cmd.Run()
}
