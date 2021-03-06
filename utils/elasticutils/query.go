package elasticutils

import (
	"context"

	"github.com/pkg/errors"

	"gopkg.in/olivere/elastic.v6"
)

// CreateBulk creates a lazily computed set of objects in bulk
func CreateBulk(client *elastic.Client, index string, tp string, objects chan interface{}) (err error) {
	// create a new bulk request
	bulkRequest := client.Bulk()

	for object := range objects {
		req := elastic.NewBulkIndexRequest().Index(index).Type(tp).Doc(object)
		bulkRequest.Add(req)
	}

	ctx := context.Background()
	res, err := bulkRequest.Do(ctx)
	err = errors.Wrap(err, "bulkRequest.Do failed")

	if err == nil && res.Errors {
		err = errors.New("[CreateBulk] Elasticsearch reported Errors=true")
	}

	return
}

// UpdateAll updates all objects inside a given index
func UpdateAll(client *elastic.Client, index string, tp string, script *elastic.Script) (err error) {
	ctx := context.Background()
	res, err := client.UpdateByQuery(index).Type(tp).Query(elastic.NewMatchAllQuery()).Script(script).Do(ctx)
	err = errors.Wrap(err, "client.UpdateByQuery failed")

	if err == nil && res.TimedOut {
		err = errors.New("[UpdateAll] Elasticsearch reported TimedOut=true")
	}

	return
}

// DeleteBulk deletes objects subject to a given query
func DeleteBulk(client *elastic.Client, index string, tp string, query elastic.Query) (err error) {
	ctx := context.Background()
	res, err := client.DeleteByQuery(index).Type(tp).Query(query).Do(ctx)
	err = errors.Wrap(err, "client.DeleteByQuery failed")

	if err == nil && res.TimedOut {
		err = errors.New("[DeleteBulk] Elasticsearch reported TimedOut=true")
	}

	return
}

// Count counts all objects subject to a query
func Count(client *elastic.Client, index string, tp string, query elastic.Query) (int64, error) {
	ctx := context.Background()
	return client.Count(index).Type(tp).Query(query).Do(ctx)
}
