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
	MWSPageSize  int64 // page size for MathWebSearch pages
	TemaPageSize int64 // page zie for TemaSearch pages

	PoolSize int64 // maximum number of parallel queries
}

// NewConnection makes a new TemaSearch Connection
func NewConnection(mws *mws.Connection, tema *tema.Connection) *Connection {
	return &Connection{
		config: &Configuration{
			MWSPageSize:  1,
			TemaPageSize: 10,
			PoolSize:     10,
		},
		mws:  mws,
		tema: tema,
	}
}
