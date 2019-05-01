package connection

import (
	"time"

	elastic "gopkg.in/olivere/elastic.v6"
)

// ElasticConnection represents a connection to an ElasticSearch instance
type ElasticConnection struct {
	port     int    // port number used
	hostname string // hostnameused

	Client *elastic.Client       // underlying http client
	Config *ElasticConfiguration // tema configuration
}

// ElasticConfiguration represents a TemaSearch Configuration
type ElasticConfiguration struct {
	HarvestIndex string
	HarvestType  string

	SegmentIndex string
	SegmentType  string

	Timeout time.Duration

	PoolSize    int
	MaxPageSize int64
}

// NewElasticConnection initializes a new Tema connection
func NewElasticConnection(port int, hostname string) (conn *ElasticConnection, err error) {
	conn = &ElasticConnection{
		port:     port,
		hostname: hostname,

		Config: &ElasticConfiguration{
			HarvestIndex: "tema",
			HarvestType:  "_doc",

			SegmentIndex: "tema-segments",
			SegmentType:  "_doc",

			Timeout: 5 * time.Second,

			PoolSize:    10,
			MaxPageSize: 10,
		},
	}

	// and validate the connection
	err = Validate(conn.port, conn.hostname)
	return
}

// connect connects to this connection
func (conn *ElasticConnection) connect() (err error) {
	// create a new elasticsearch server
	client, err := elastic.NewClient(elastic.SetURL(MakeURL(conn.port, conn.hostname, "")), elastic.SetSniff(false), elastic.SetHealthcheckTimeoutStartup(conn.Config.Timeout))
	if err != nil {
		return
	}

	conn.Client = client
	return
}

// Close closes this connection
func (conn *ElasticConnection) Close() error {
	if conn.Client != nil {
		conn.Client.Stop()
		conn.Client = nil
	}
	return nil
}

func init() {
	// ensure at compile time that ElasticConnection implements Connection
	var _ Connection = (*ElasticConnection)(nil)
}
