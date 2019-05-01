package connection

import (
	"time"

	"gopkg.in/olivere/elastic.v6"
)

// TemaConnection represents a connection to TemaSearch
type TemaConnection struct {
	port     int    // port number used
	hostname string // hostnameused

	Client *elastic.Client    // underlying http client
	Config *TemaConfiguration // tema configuration
}

// TemaConfiguration represents a TemaSearch Configuration
type TemaConfiguration struct {
	HarvestIndex string
	HarvestType  string

	SegmentIndex string
	SegmentType  string

	Timeout time.Duration

	PoolSize    int
	MaxPageSize int64
}

// NewTemaConnection initializes a new Tema connection
func NewTemaConnection(port int, hostname string) (conn *TemaConnection, err error) {
	conn = &TemaConnection{
		port:     port,
		hostname: hostname,

		Config: &TemaConfiguration{
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
func (conn *TemaConnection) connect() (err error) {
	// create a new elasticsearch server
	client, err := elastic.NewClient(elastic.SetURL(MakeURL(conn.port, conn.hostname, "")), elastic.SetSniff(false), elastic.SetHealthcheckTimeoutStartup(conn.Config.Timeout))
	if err != nil {
		return
	}

	conn.Client = client
	return
}

// close closes this connection
func (conn *TemaConnection) close() {
	if conn.Client != nil {
		conn.Client.Stop()
		conn.Client = nil
	}
}

func init() {
	// ensure at compile time that TemaConnection implements Connection
	var _ Connection = (*TemaConnection)(nil)
}
