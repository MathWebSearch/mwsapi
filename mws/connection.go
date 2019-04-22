package mws

import (
	"fmt"
	"net/http"
)

// Connection represents a connection to a MathWebSearch Daemon
type Connection struct {
	client   *http.Client
	Hostname string
	Port     int
}

// NewConnection creates a new connection
func NewConnection(Host string, Port int) *Connection {
	return &Connection{
		client:   &http.Client{},
		Hostname: Host,
		Port:     Port,
	}
}

// URL returns the URL of the connection
func (conn *Connection) URL() string {
	return fmt.Sprintf("http://%s:%d/", conn.Hostname, conn.Port)
}
