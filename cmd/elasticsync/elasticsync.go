package main

import (
	"fmt"
	"os"

	"github.com/MathWebSearch/mwsapi/cmd/elasticsync/args"
	"github.com/MathWebSearch/mwsapi/tema"
	"github.com/MathWebSearch/mwsapi/tema/sync"
)

func main() {
	// parse and validate arguments
	a := args.ParseArgs(os.Args)
	if !a.Validate() {
		os.Exit(1)
		return
	}

	// connect to elasticsearch
	fmt.Printf("Connecting to %q ...\n", a.ElasticURL())
	connection := tema.Connect(a.ElasticHost, a.ElasticPort)
	fmt.Println("Connected. ")

	// make a sync process
	process := sync.NewProcess(connection, a.IndexDir, a.Quiet)
	stats, err := process.Run()

	// if there was an error, print it
	if err != nil {
		fmt.Printf(err.Error())
		os.Exit(1)
	}

	// print the stats, and then exit
	fmt.Printf("Segments: %s. \n", stats.String())
	fmt.Println("Finished, exiting. ")
}

func die(err error) {

	if err != nil {
		panic(err)
	} else {
		panic("Something went wrong")
	}
}

// TODO: 1. Wait for tema-search to be up on the given port
// 2. Check if we have to run setup
// 3. Hash the directory; if it has changed clear out and fully re-index
