package elasticutils

import (
	"time"

	"gopkg.in/olivere/elastic.v6"
)

// TryConnect connects to the server
func TryConnect(retryInterval time.Duration, maxRetries int, onRetry func(error), funcs ...elastic.ClientOptionFunc) (cli *elastic.Client, err error) {
	counter := 0
	for {
		cli, err = elastic.NewClient(funcs...)

		if err == nil {
			break
		}

		// if the counter has been reached, break
		counter++
		if maxRetries >= 0 && counter > maxRetries {
			break
		}

		// and wait for next time
		if onRetry != nil {
			onRetry(err)
		}
		time.Sleep(retryInterval)

	}
	return
}
