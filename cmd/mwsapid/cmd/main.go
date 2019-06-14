package cmd

import (
	"log"

	"github.com/MathWebSearch/mwsapi/engine"
	"github.com/MathWebSearch/mwsapi/engine/mwsengine"
	"github.com/MathWebSearch/mwsapi/engine/temaengine"
	"github.com/pkg/errors"
)

// Main represents the main interface of the mwsquery command
func Main(a *Args) {
	// make a new connection to mws
	mws := &mwsengine.MWSHandler{
		Host: a.MWSHost, Port: a.MWSPort,
	}
	err := mws.Connect()
	if err != nil {
		err = errors.Wrap(err, "unable to connect to MWS server")
		log.Fatalf("%v", err)
		return
	}

	// make a new connection to tema
	tema := &temaengine.TemaHandler{
		MWSPort: a.MWSPort, MWSHost: a.MWSHost,
		ElasticPort: a.ElasticPort, ElasticHost: a.ElasticHost,
	}
	err = tema.Connect()
	if err != nil {
		err = errors.Wrap(err, "unable to connect to Temasearch server")
		log.Fatalf("%v", err)
		return
	}

	// create the server
	server := engine.NewServer()
	server.AddHandler(mws)
	server.AddHandler(tema)

	// and start servering
	log.Printf("Listening on %s:%d", a.Host, a.Port)
	log.Fatal(server.ListenAndServe(a.Host, a.Port))
}
