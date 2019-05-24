package elasticutils

import (
	"context"

	"github.com/pkg/errors"
	"gopkg.in/olivere/elastic.v6"
)

// CreateIndex creates an index unless it already exists
func CreateIndex(client *elastic.Client, index string, mapping interface{}) (created bool, err error) {
	ctx := context.Background()

	// check if the index exists
	exists, err := client.IndexExists(index).Do(ctx)
	err = errors.Wrap(err, "client.IndexExists failed")
	if err != nil {
		return
	}

	// create it if not
	if !exists {
		res, err := client.CreateIndex(index).BodyJson(mapping).Do(ctx)
		err = errors.Wrap(err, "client.CreateIndex failed")
		if err == nil && !res.Acknowledged {
			err = errors.New("[CreateIndex] Elasticsearch reported acknowledged=false")
		}

		if err != nil {
			return false, err
		}
		created = true
	}

	return
}

// RefreshIndex refreshes an index
func RefreshIndex(client *elastic.Client, indices ...string) (err error) {
	ctx := context.Background()
	res, err := client.Refresh(indices...).Do(ctx)
	err = errors.Wrap(err, "client.Refresh failed")

	if err == nil && res.Shards.Successful <= 0 {
		err = errors.New("[RereshIndex] Elasticsearch reported 0 successful shards")
	}

	return
}

// FlushIndex flushes an index
func FlushIndex(client *elastic.Client, indices ...string) (err error) {
	ctx := context.Background()
	res, err := client.Flush(indices...).Do(ctx)
	err = errors.Wrap(err, "client.Flush failed")

	if err == nil && res.Shards.Successful <= 0 {
		err = errors.New("[Flush] Elasticsearch reported 0 successful shards")
	}

	return
}
