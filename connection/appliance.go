package connection

import (
	"github.com/MathWebSearch/mwsapi/utils/gogroup"
)

// ApplianceConnection represents a connection to a MathWebSearch appliance
type ApplianceConnection struct {
	MWS  *MWSConnection  // connection to MathWebSearch
	Tema *TemaConnection // connection to Tema
}

// NewApplianceConnection makes a new connection to an appliance
func NewApplianceConnection(MWSPort int, MWSHost string, TemaPort int, TemaHost string) (conn *ApplianceConnection, err error) {
	conn = &ApplianceConnection{}

	// create the MWS Connection
	conn.MWS, err = NewMWSConnection(MWSPort, MWSHost)
	if err != nil {
		return
	}

	// create the tema connection
	conn.Tema, err = NewTemaConnection(TemaPort, TemaHost)
	return
}

func (conn *ApplianceConnection) connect() (err error) {
	group := gogroup.NewWorkGroup(-1, false)

	// connect to mws
	mws := gogroup.GroupJob(func(_ func(func())) error {
		return conn.MWS.connect()
	})
	group.Add(&mws)

	// connect to tema
	tema := gogroup.GroupJob(func(_ func(func())) error {
		return conn.Tema.connect()
	})
	group.Add(&tema)

	// wait for both to finish
	err = group.Wait()

	// if either of the connections failed
	// then we need to disconnect the other
	if err != nil {
		conn.close()
	}

	return
}

func (conn *ApplianceConnection) close() {
	// close mws
	if conn.MWS != nil {
		conn.MWS.close()
	}

	// close tema
	if conn.Tema != nil {
		conn.Tema.close()
	}
}

func init() {
	// ensure at compile time that ApplianceConnection implements Connection
	var _ Connection = (*ApplianceConnection)(nil)
}
