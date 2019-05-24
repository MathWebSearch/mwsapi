package connection

import (
	"github.com/MathWebSearch/mwsapi/utils/gogroup"
	"github.com/pkg/errors"
)

// TemaConnection represents a connection to a TemaSearch instance, that is a joined (MathWebSearch, ElasticSearch) instance
type TemaConnection struct {
	MWS     *MWSConnection     // connection to MathWebSearch
	Elastic *ElasticConnection // connection to ElasticSearch
}

// NewTemaConnection makes a new connection to TemaSearch
func NewTemaConnection(MWSPort int, MWSHost string, ElasticPort int, ElasticHost string) (conn *TemaConnection, err error) {
	conn = &TemaConnection{}

	// create the MWS Connection
	conn.MWS, err = NewMWSConnection(MWSPort, MWSHost)
	err = errors.Wrap(err, "NewMWSConnection failed")
	if err != nil {
		return
	}

	// create the tema connection
	conn.Elastic, err = NewElasticConnection(ElasticPort, ElasticHost)
	err = errors.Wrap(err, "NewElasticConnection failed")
	return
}

func (conn *TemaConnection) connect() (err error) {
	group := gogroup.NewWorkGroup(-1, false)

	// connect to mws
	mws := gogroup.GroupJob(func(_ func(func())) error {
		return errors.Wrap(conn.MWS.connect(), "conn.MWS.connect failed")
	})
	group.Add(&mws)

	// connect to tema
	tema := gogroup.GroupJob(func(_ func(func())) error {
		return errors.Wrap(conn.Elastic.connect(), "conn.Elastic.connect failed")
	})
	group.Add(&tema)

	// wait for both to finish
	err = group.Wait()
	err = errors.Wrap(err, "group.Wait failed")

	// if either of the connections failed
	// then we need to disconnect the other
	if err != nil {
		conn.Close()
	}

	return
}

// Close closes this connection
func (conn *TemaConnection) Close() error {
	// close mws
	if conn.MWS != nil {
		conn.MWS.Close()
	}

	// close tema
	if conn.Elastic != nil {
		conn.Elastic.Close()
	}

	return nil
}

func init() {
	// ensure at compile time that TemaConnection implements Connection
	var _ Connection = (*TemaConnection)(nil)
}
