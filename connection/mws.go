package connection

import (
	"net"
	"net/http"
	"strconv"
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
	// ping the connection
	err = conn.ping()
	if err != nil {
		return
	}

	conn.Client = &http.Client{}
	return
}

func (conn *MWSConnection) ping() (err error) {
	// connect via tcp
	c, err := net.DialTimeout("tcp", net.JoinHostPort(conn.hostname, strconv.Itoa(conn.port)), conn.Config.Timeout)
	if err != nil {
		return
	}

	// and close the connection immediatly
	c.Close()
	return
}

// closes the connection to MathWebSearch
func (conn *MWSConnection) close() {
	conn.Client = nil
}

func init() {
	// ensure at compile time that MWSConnection implements Connection
	var _ Connection = (*MWSConnection)(nil)
}
