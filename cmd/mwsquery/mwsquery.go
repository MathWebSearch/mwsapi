package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/MathWebSearch/mwsapi/cmd/mwsquery/args"
	"github.com/MathWebSearch/mwsapi/mws"
)

func main() {
	// parse and validate arguments
	a := args.ParseArgs(os.Args)
	if !a.Validate() {
		os.Exit(1)
		return
	}

	connection := mws.NewConnection(a.MWSHost, a.MWSPort)
	query := &mws.Query{
		Expressions: a.Expressions,
		MwsIdsOnly:  a.MWSIdsOnly,
	}

	var res interface{}
	var err error

	if a.Count {
		res, err = mws.CountQuery(connection, query)
	} else {
		res, err = mws.RunQuery(connection, query, a.From, a.Size)
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
