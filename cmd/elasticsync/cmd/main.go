package cmd

import (
	"time"

	"github.com/MathWebSearch/mwsapi/connection"
	"github.com/MathWebSearch/mwsapi/engine/elasticsync"
)

// Main represents the main interface of the elasticsync command
func Main(a *Args) (res interface{}, err error) {
	// make a new connection
	c, err := connection.NewElasticConnection(a.ElasticPort, a.ElasticHost)
	if err != nil {
		return
	}

	// make a new process
	process := elasticsync.NewProcess(c, a.IndexDir, a.Quiet, a.Force)

	// connect to elasticsearch
	process.Printf(nil, "Connecting to %q ...\n", a.ElasticURL())
	err = connection.AwaitConnect(c, 5*time.Second, -1, func(e error) {
		process.Printf(nil, "  Connection failed: %q, trying again in 5 seconds.\n", e.Error())
	})
	process.PrintlnOK(nil, "Connected. ")
	defer c.Close()

	{
		// run the sync process
		stats, err := process.Run()
		if err == nil && a.Normalize && stats != nil {
			stats.Normalize()
		}

		// and return
		return stats, err
	}
}
