package main

import (
	"os"
	"time"

	"github.com/MathWebSearch/mwsapi/cmd/elasticsync/args"
	"github.com/MathWebSearch/mwsapi/connection"
	"github.com/MathWebSearch/mwsapi/engine/elasticsync"
	"github.com/MathWebSearch/mwsapi/utils"
)

func main() {
	// parse and validate arguments
	a := args.ParseArgs(os.Args)
	if !a.Validate() {
		os.Exit(1)
		return
	}

	// make a new connection
	c, err := connection.NewTemaConnection(a.ElasticPort, a.ElasticHost)
	if err != nil {
		panic(err)
	}

	// make a new process
	process := elasticsync.NewProcess(c, a.IndexDir, a.Quiet, a.Force)

	// connect to elasticsearch
	process.Printf(nil, "Connecting to %q ...\n", a.ElasticURL())
	err = connection.AwaitConnect(c, 5*time.Second, -1, func(e error) {
		process.Printf(nil, "  Connection failed: %q, trying again in 5 seconds.\n", e.Error())
	})
	process.PrintlnOK(nil, "Connected. ")

	// make a sync process
	stats, err := process.Run()
	if err == nil && a.Normalize {
		stats.Normalize()
	}

	// and output
	utils.OutputJSONOrErr(stats, err)
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
