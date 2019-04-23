package temasearch

import (
	"github.com/MathWebSearch/mwsapi/mws"
	"github.com/MathWebSearch/mwsapi/tema"
)

// Connection represents a unified TemaSearch Connection
type Connection struct {
	config *Configuration

	// connections to the daemons
	mws  *mws.Connection
	tema *tema.Connection
}

// Configuration represents a Configuration of a TemaSearch Connection
type Configuration struct {
	MWSPageSize  int64
	TemaPageSize int64
}

// NewConnection makes a new TemaSearch Connection
func NewConnection(mws *mws.Connection, tema *tema.Connection) *Connection {
	return &Connection{
		config: &Configuration{
			MWSPageSize:  1,
			TemaPageSize: 10,
		},
		mws:  mws,
		tema: tema,
	}
}
