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

	// connect to mws + temasearch
	mwsconnection := mws.NewConnection(a.MWSHost, a.MWSPort)
	temaconnection, err := tema.Connect(a.ElasticHost, a.ElasticPort)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	// make the final connection
	connection := temasearch.NewConnection(mwsconnection, temaconnection)

	// build the query
	query := &temasearch.Query{
		Expressions: a.Expressions,
		Text:        a.Text,
	}

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
