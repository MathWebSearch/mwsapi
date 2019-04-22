package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/MathWebSearch/mwsapi/cmd/temaquery/args"
	"github.com/MathWebSearch/mwsapi/tema/query"

	"github.com/MathWebSearch/mwsapi/tema"
)

func main() {
	// parse and validate arguments
	a := args.ParseArgs(os.Args)
	if !a.Validate() {
		os.Exit(1)
		return
	}

	// connect to tema-search
	connection, err := tema.Connect(a.ElasticHost, a.ElasticPort)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	var res interface{}

	if a.Count {
		res, err = RunCountQuery(connection, a)
	} else if a.DocumentPhaseOnly {
		res, err = RunDocumentQuery(connection, a)
	} else {
		res, err = RunBothPhases(connection, a)
	}

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	// stdout the json
	bytes, _ := json.MarshalIndent(res, "", "  ")
	if err != nil {
		fmt.Printf("%#v\n", err)
		os.Exit(1)
	}
	os.Stdout.Write(bytes)
	fmt.Println("")

}

// RunDocumentQuery runs only the document query
func RunDocumentQuery(connection *tema.Connection, a *args.Args) (res interface{}, err error) {
	res, err = query.RunDocumentQuery(connection, &query.Query{
		Text:             a.Text,
		MathWebSearchIDs: a.IDs,
	}, a.From, a.Size)
	return
}

// RunBothPhases runs both phases
func RunBothPhases(connection *tema.Connection, a *args.Args) (res interface{}, err error) {
	res, err = query.RunQuery(connection, &query.Query{
		Text:             a.Text,
		MathWebSearchIDs: a.IDs,
	}, a.From, a.Size)
	return
}

// RunCountQuery runs a count query
func RunCountQuery(connection *tema.Connection, a *args.Args) (res interface{}, err error) {
	res, err = query.CountQuery(connection, &query.Query{
		Text:             a.Text,
		MathWebSearchIDs: a.IDs,
	})
	return
}
