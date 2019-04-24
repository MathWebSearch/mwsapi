package tema

import (
	"fmt"
	"time"

	"github.com/MathWebSearch/mwsapi/elasticutils"
	"gopkg.in/olivere/elastic.v6"
)

// Connection represents a connection to a TemaSearch Instance
type Connection struct {
	Config *Configuration
	Client *elastic.Client
}

// Configuration represents a TemaSearch Configuration
type Configuration struct {
	HarvestIndex string
	HarvestType  string

	SegmentIndex string
	SegmentType  string

	PoolSize int
}

type j map[string]interface{}

// ConnectionFromClient creates a new connection from an existing elastic client
func ConnectionFromClient(client *elastic.Client) *Connection {
	return &Connection{
		Config: &Configuration{
			HarvestIndex: "tema",
			HarvestType:  "_doc",

			SegmentIndex: "tema-segments",
			SegmentType:  "_doc",

			PoolSize: 10,
		},
		Client: client,
	}
}

// WaitConnect repeatetly connects to a TemaSearch instance until the connection succeeds
func WaitConnect(Host string, Port int) *Connection {
	client, _ := elasticutils.TryConnect(5*time.Second, -1, func(e error) {
		fmt.Printf("  Unable to connect: %s. Trying again in 5 seconds. \n", e)
	}, elastic.SetURL(fmt.Sprintf("http://%s:%d", Host, Port)), elastic.SetSniff(false))

	return ConnectionFromClient(client)
}

// Connect connects to a TemaSearch instance or returns an error
func Connect(Host string, Port int) (conn *Connection, err error) {
	cli, err := elasticutils.TryConnect(0, 0, nil, elastic.SetSniff(false))
	if err != nil {
		return
	}

	conn = ConnectionFromClient(cli)
	return
}
