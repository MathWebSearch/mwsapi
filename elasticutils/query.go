package elasticutils

import (
	"context"
	"errors"

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

	if err == nil && res.Errors {
		err = errors.New("Bulk() reported Errors=true")
	}

	return
}

// UpdateAll updates all objects inside a given index
func UpdateAll(client *elastic.Client, index string, tp string, script *elastic.Script) (err error) {
	ctx := context.Background()
	res, err := client.UpdateByQuery(index).Type(tp).Query(elastic.NewMatchAllQuery()).Script(script).Do(ctx)

	if err == nil && res.TimedOut {
		err = errors.New("UpdateByQuery() reported TimedOut=true")
	}

	return
}

// DeleteBulk deletes objects subject to a given query
func DeleteBulk(client *elastic.Client, index string, tp string, query elastic.Query) (err error) {
	ctx := context.Background()
	res, err := client.DeleteByQuery(index).Type(tp).Query(query).Do(ctx)

	if err == nil && res.TimedOut {
		err = errors.New("DeleteByQuery() reported TimedOut=true")
	}

	return
}
