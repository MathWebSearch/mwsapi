package elasticutils

import (
	"time"

	"gopkg.in/olivere/elastic.v6"
)

// Connect connects to the server
func Connect(retryInterval time.Duration, onRetry func(error), funcs ...elastic.ClientOptionFunc) (cli *elastic.Client) {
	var err error
	for {
		cli, err = elastic.NewClient(funcs...)

		if err == nil {
			break
		}

		// and wait for next time
		onRetry(err)
		time.Sleep(retryInterval)
	}
	return
}
