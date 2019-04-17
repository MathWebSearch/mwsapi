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
}

type j map[string]interface{}

// Connect connects to a temasearch instance until it is connected
func Connect(Host string, Port int) *Connection {
	client := elasticutils.Connect(5*time.Second, func(e error) {
		fmt.Printf("  Unable to connect: %s. Trying again in 5 seconds. \n", e)
	}, elastic.SetURL(fmt.Sprintf("http://%s:%d", Host, Port)), elastic.SetSniff(false))

	return &Connection{
		Config: &Configuration{
			HarvestIndex: "tema",
			HarvestType:  "_doc",

			SegmentIndex: "tema-segments",
			SegmentType:  "_doc",
		},
		Client: client,
	}
}
