package integrationtest

import (
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"testing"
	"time"

	"github.com/MathWebSearch/mwsapi/utils/gogroup"
	"github.com/pkg/errors"
)

// StartDockerService starts docker-compose configured services from the given servicefile
// and then blocks until a url returns http status code 200
func StartDockerService(service string, urls ...string) (client *http.Client, err error) {
	err = runExternalCommand("docker-compose", "-f", service, "up", "--force-recreate", "-d")
	err = errors.Wrap(err, "runExternalCommand failed")
	if err != nil {
		return
	}

	// create a new client and wait group
	client = &http.Client{}

	group := gogroup.NewWorkGroup(0, true)

	for _, url := range urls {
		func(url string) {
			job := gogroup.GroupJob(func(sync func(func())) error {
				if testing.Verbose() {
					sync(func() {
						fmt.Printf("Waiting for %q to return HTTP 200 ", url)
					})
				}

				for {
					if testing.Verbose() {
						sync(func() {
							fmt.Print(".")
						})
					}
					res, err := client.Get(url)
					err = errors.Wrap(err, "client.Get failed")
					if err == nil && res.StatusCode == 200 {
						break
					}
					time.Sleep(1 * time.Second) // wait for the next one
				}

				if testing.Verbose() {
					sync(func() { fmt.Println(" ok") })
				}

				return nil
			})

			group.Add(&job)
		}(url)
	}

	// wait for all of them to be done
	group.Wait()
	return
}

// StopDockerService gracefully stops and removes the docker-compose configured services from servicefile
// including all volumes
func StopDockerService(servicefile string) (err error) {
	// stop the service
	err = runExternalCommand("docker-compose", "-f", servicefile, "down", "-v")
	err = errors.Wrap(err, "runExternalCommand failed")
	return
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
	err = cmd.Run()
	err = errors.Wrap(err, "cmd.Run failed")
	return
}
