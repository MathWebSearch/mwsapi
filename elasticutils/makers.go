package elasticutils

import (
	"context"
	"errors"
	"io"

	"gopkg.in/olivere/elastic.v6"
)

// NewObjectFromID fetches a new EC object from the server
func NewObjectFromID(client *elastic.Client, index string, tp string, id string) (obj *Object, err error) {
	// create an empty object
	obj = &Object{client, index, tp, id, nil, nil}

	// reload it from the db, clear if it fails
	err = obj.Reload()
	if err != nil {
		obj = nil
	}
	return
}

// NewObjectFromFields creates a new ec object on the server
func NewObjectFromFields(client *elastic.Client, index string, tp string, Data interface{}) (obj *Object, err error) {
	obj = &Object{client, index, tp, "", nil, nil}

	// pack the fields into the object
	err = obj.Pack(Data)
	if err != nil {
		return
	}

	err = obj.Index()
	return
}

// NewObjectFromHit creates a new object using a SearchHit
func NewObjectFromHit(client *elastic.Client, hit *elastic.SearchHit) (obj *Object, err error) {
	// make a new object with index and types
	obj = &Object{client, hit.Index, hit.Type, hit.Id, nil, nil}

	//update it from the hit
	err = obj.updateFromHit(hit)
	if err != nil {
		obj = nil
	}

	return
}

// FetchObjects fetches objects subject to an exact query
func FetchObjects(client *elastic.Client, index string, tp string, query elastic.Query) <-chan *Object {

	ctx := context.Background()
	scroll := client.Scroll(index).Type(tp).Query(query)

	hits := make(chan *Object)

	go func() {
		defer close(hits)

		for {
			results, err := scroll.Do(ctx)
			if err == io.EOF {
				return
			}
			if err != nil {
				return
			}

			for _, hit := range results.Hits.Hits {
				obj, err := NewObjectFromHit(client, hit)
				if err == nil {
					hits <- obj
				}
			}
		}

	}()

	return hits
}

// ResultsPage represents a single page of results
type ResultsPage struct {
	Hits  []*Object
	Total int64

	From int64
	Size int64
}

// FetchObjectsPage fetches a set of objects given a specific slicing
func FetchObjectsPage(client *elastic.Client, index string, tp string, query elastic.Query, highlight *elastic.Highlight, from int64, size int64) (page *ResultsPage, err error) {
	ctx := context.Background()
	search := client.Search(index).Type(tp).Query(query).From(int(from)).Size(int(size))
	if highlight != nil {
		search = search.Highlight(highlight)
	}

	results, err := search.Do(ctx)

	if err == nil && results.TimedOut {
		err = errors.New("Search() reported TimedOut=true")
	}

	if err != nil {
		return
	}

	page = &ResultsPage{
		From: from,
		Size: size,
		Hits: make([]*Object, len(results.Hits.Hits)),
	}

	// iterate over the hits
	for i, hit := range results.Hits.Hits {
		obj, err := NewObjectFromHit(client, hit)
		if err != nil {
			return nil, err
		}

		page.Hits[i] = obj
	}

	// count the hits
	page.Total = results.TotalHits()

	// and return
	return
}

// FetchObject fetches a single object from the database or returns nil
func FetchObject(client *elastic.Client, index string, tp string, query elastic.Query, highlight *elastic.Highlight) (obj *Object, err error) {

	// fetch a page of objects
	results, err := FetchObjectsPage(client, index, tp, query, highlight, 0, 1)
	if err != nil {
		return
	}

	// no hits => no returned result
	if results.Total == 0 {
		return
	}

	// and return the one (and only) result
	obj = results.Hits[0]
	return
}

// FetchOrCreateObject fetches the object returned from the query, or creates a new one if no result is retrieved
func FetchOrCreateObject(client *elastic.Client, index string, tp string, query elastic.Query, Data interface{}) (obj *Object, created bool, err error) {
	// first try and fetch the object
	obj, err = FetchObject(client, index, tp, query, nil)
	if err != nil || obj != nil {
		return
	}

	// if that fails create it
	obj, err = NewObjectFromFields(client, index, tp, Data)
	if err != nil {
		created = true
	}

	return
}
