package main

import (
	"fmt"

	"github.com/MathWebSearch/mwsapi/tema/query"

	"github.com/MathWebSearch/mwsapi/tema"
)

func main() {
	// connect to tema-search
	connection := tema.Connect("0.0.0.0", 9200)

	res, err := query.RunDocumentQuery(connection, &query.Query{
		Text: "math",
	}, 0, 10)
	fmt.Printf("res := %#v\n", res)
	fmt.Printf("err := %#v\n", err)
}
