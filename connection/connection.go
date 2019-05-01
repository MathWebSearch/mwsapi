package connection

import (
	"fmt"
	"time"
)

// Connection represents the connection to a given daemon
type Connection interface {
	connect() error // Attempt to connect to the server
	close()         // Close the connection (if open)
}

// Validate validates details to a connection
func Validate(port int, hostname string) (e error) {
	if port <= 0 || port > 65535 {
		return fmt.Errorf("%d is not a valid port number. Port numbers are in the range (1, 65535) inclusive. ", port)
	}

	if hostname == "" {
		return fmt.Errorf("%q is not a valid hostname. It should not be empty. ", hostname)
	}

	return
}

// MakeURL makes the url to a connection
// protocol defaults to "http" when empty
func MakeURL(port int, hostname string, protocol string) string {
	if protocol == "" {
		protocol = "http"
	}
	return fmt.Sprintf("%s://%s:%d", protocol, hostname, port)
}

// Connect connects to the daemon represented by this connection
func Connect(c Connection) error {
	return c.connect()
}

// AwaitConnect repeatedly tries to connect to the daemon until the connection suceeds
// retryInterval is the time to wait between each successive connection attempt.
// maxRetries is the maximal number of attempts to retry the connection; a value < 0 represents an infinite attempt of attempts
// onFailure is called whenever a connection attempt fails; might be nil to not be called
func AwaitConnect(c Connection, retryInterval time.Duration, maxRetries int, onFailure func(error)) (err error) {
	for {
		// try to connect and exit of successfull
		err = c.connect()
		if err == nil {
			break
		}

		// call the failure handler
		if onFailure != nil {
			onFailure(err)
		}

		// if we have no more retries left, break
		if maxRetries == 0 {
			break
		}
		maxRetries--

		// wait for the next attempt
		time.Sleep(retryInterval)
	}
	return
}
