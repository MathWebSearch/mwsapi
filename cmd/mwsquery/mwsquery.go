package main

import (
	"fmt"

	"github.com/MathWebSearch/mwsapi/mws"
)

func main() {
	thing := &mws.Query{}
	bytes, _ := thing.ToXML()
	fmt.Println(string(bytes))
}
