package main

import (
	"fmt"

	"github.com/MathWebSearch/mwsapi/mws"
)

func main() {
	connection := mws.NewConnection("localhost", 8080)
	query := &mws.Query{
		Expressions: []string{"<mws:qvar>x</mws:qvar>"},
		MwsIdsOnly:  false,
	}

	res, err := mws.RunQuery(connection, query, 0, 10)

	if res != nil {
		fmt.Printf("res = %#v\n", res)
	}
	fmt.Printf("err = %#v\n", err)
	if err != nil {
		fmt.Println(err.Error())
	}
}
