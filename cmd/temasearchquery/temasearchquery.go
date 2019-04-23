package main

import (
	"fmt"
	"os"

	"github.com/MathWebSearch/mwsapi/cmd/temasearchquery/args"
	"github.com/MathWebSearch/mwsapi/mws"
	"github.com/MathWebSearch/mwsapi/tema"
	"github.com/MathWebSearch/mwsapi/temasearch"
	"github.com/MathWebSearch/mwsapi/utils"
)

func main() {
	// parse and validate arguments
	a := args.ParseArgs(os.Args)
	if !a.Validate() {
		os.Exit(1)
		return
	}

	// build the query
	query := &temasearch.Query{
		Expressions: a.Expressions,
		Text:        a.Text,
	}

	// setup mws connection (if needed)
	var mwsConnection *mws.Connection
	if query.NeedsMWS() {
		mwsConnection = mws.NewConnection(a.MWSHost, a.MWSPort)
	}

	// setup tema connection (if needed)
	var temaConnection *tema.Connection
	var err error
	if query.NeedsElastic() {
		// connect to mws + temasearch
		temaConnection, err = tema.Connect(a.ElasticHost, a.ElasticPort)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
	}

	// make the final connection
	connection := temasearch.NewConnection(mwsConnection, temaConnection)

	// and run it
	var res interface{}
	if a.Count {
		res, err = temasearch.CountQuery(connection, query)
	} else {
		res, err = temasearch.RunQuery(connection, query, a.From, a.Size)
	}

	// and output
	utils.OutputJSONOrErr(&res, err)
}
