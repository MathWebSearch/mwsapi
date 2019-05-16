package connection

import (
	"errors"
	"net/http"
	"time"
)

// MWSConnection represents a connection to MathWebSearch
type MWSConnection struct {
	port     int
	hostname string

	Client *http.Client
	Config *MWSConfiguration
}

// MWSConfiguration represents a MathWebSearch Configuration
type MWSConfiguration struct {
	Timeout time.Duration

	PoolSize    int
	MaxPageSize int64
}

// NewMWSConnection initializes a new MathWebSearch connection
func NewMWSConnection(port int, hostname string) (conn *MWSConnection, err error) {
	conn = &MWSConnection{
		port:     port,
		hostname: hostname,

		Config: &MWSConfiguration{
			Timeout: 5 * time.Second,

			PoolSize:    10,
			MaxPageSize: 10,
		},
	}

	// and validate the connection
	err = Validate(conn.port, conn.hostname)
	return
}

// URL returns the URL to this connection
func (conn *MWSConnection) URL() string {
	return MakeURL(conn.port, conn.hostname, "http")
}

// connect connects to this connection
func (conn *MWSConnection) connect() (err error) {
	// create a new http client
	conn.Client = &http.Client{}

	// ping and make sure the connection actually works
	err = conn.ping()
	if err != nil {
		conn.Client = nil
	}

	return
}

func (conn *MWSConnection) ping() (err error) {
	res, err := conn.Client.Get(conn.URL())
	if err != nil {
		return
	}

	// check that the status code is 200
	if res.StatusCode != 200 {
		err = errors.New("[MWSConnection.ping] MathWebSearch did not respond with code HTTP 200")
	}

	return
}

// Close closes the connection to MathWebSearch
func (conn *MWSConnection) Close() error {
	conn.Client = nil
	return nil
}

func init() {
	// ensure at compile time that MWSConnection implements Connection
	var _ Connection = (*MWSConnection)(nil)
}
