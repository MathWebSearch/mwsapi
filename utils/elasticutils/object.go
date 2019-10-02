package elasticutils

import (
	"context"
	"encoding/json"

	jsoniter "github.com/json-iterator/go"

	"github.com/pkg/errors"

	"gopkg.in/olivere/elastic.v6"
)

// Object represents an object within elasticsearch
type Object struct {

	// the client that this object uses
	client *elastic.Client

	index string // index the object resides in
	tp    string // type the object resides in
	id    string // id of this object

	Source *json.RawMessage // the source of this object (if any)
	Hit    *HitInfo         // information about this object as a result (if available)
}

// HitInfo represents information about this object as a hit
type HitInfo struct {
	Highlight *elastic.SearchHitHighlight
	Score     *float64
}

// IsIndexed checks if an object is still indexed
func (obj *Object) IsIndexed() bool {
	return obj.id != ""
}

// GetID gets the id of this object
func (obj *Object) GetID() string {
	return obj.id
}

// Index indexes this object in the database
func (obj *Object) Index() (err error) {
	// if we already have an id, it is already indexed
	if obj.IsIndexed() {
		return errors.New("[Object.Index] Already indexed")
	}

	// grab the fields
	doc, err := obj.UnpackFields()
	err = errors.Wrap(err, "obj.UnpackFields failed")
	if err != nil {
		return
	}

	// perform the Indexing operation
	ctx := context.Background()
	result, err := obj.client.Index().Index(obj.index).Type(obj.tp).BodyJson(doc).Do(ctx)
	err = errors.Wrap(err, "obj.client.Index failed")
	if err == nil && result.Shards.Successful <= 0 {
		err = errors.New("[Object.Index] Elasticsearch reported 0 successful shards")
	}

	if err != nil {
		return
	}

	// grab the new object id
	obj.id = result.Id

	return
}

// Reload reloads this object from the database
func (obj *Object) Reload() (err error) {

	if !obj.IsIndexed() {
		return errors.New("[Object.Reload] Not indexed")
	}

	ctx := context.Background()

	// fetch from the db and return unless an error occured
	result, err := obj.client.Get().Index(obj.index).Type(obj.tp).Id(obj.id).Do(ctx)
	err = errors.Wrap(err, "obj.client.Get failed")
	if err == nil && !result.Found {
		err = errors.New("[Object.Reload] Did not find object")
	}

	if err != nil {
		return
	}

	obj.Source = result.Source

	return
}

func (obj *Object) updateFromHit(source *elastic.SearchHit) error {
	obj.Hit = &HitInfo{
		Score:     source.Score,
		Highlight: &source.Highlight,
	}
	obj.Source = source.Source

	return nil
}

// Save saves this object into the database
func (obj *Object) Save() (err error) {
	// check if we are indexed
	if !obj.IsIndexed() {
		return errors.New("[Object.Save] Not indexed")
	}

	// grab the fields
	doc, err := obj.UnpackFields()
	err = errors.Wrap(err, "obj.UnpackFields failed")
	if err != nil {
		return
	}

	// and replace it in the database
	ctx := context.Background()
	res, err := obj.client.Update().Index(obj.index).Type(obj.tp).Id(obj.id).Doc(doc).Do(ctx)
	err = errors.Wrap(err, "obj.client.Update failed")

	//check errors
	if err == nil && (res.Result != "noop" && res.Shards.Successful <= 0) {
		err = errors.New("[Object.Save] Elasticsearch reported non-noop result with 0 successful shards ")
	}

	return
}

// Delete deletes this object
func (obj *Object) Delete() (err error) {
	if !obj.IsIndexed() {
		return errors.New("[Object.Delete] Not indexed")
	}

	// just clears the object
	ctx := context.Background()
	res, err := obj.client.Delete().Index(obj.index).Type(obj.tp).Id(obj.id).Do(ctx)
	err = errors.Wrap(err, "obj.client.Delete failed")

	if err == nil && res.Result != "deleted" {
		err = errors.New("[Object.Delete] Elasticsearch did not report deleted result")
	}

	// id no longer valid
	if err == nil {
		obj.id = ""
	}

	return
}

// Unpack will unpack this object as json
func (obj *Object) Unpack(v interface{}) (err error) {
	err = jsoniter.Unmarshal(*obj.Source, v)
	err = errors.Wrap(err, "jsoniter.Unmarshal failed")
	return
}

// UnpackFields unpacks this object into a set of fields
func (obj *Object) UnpackFields() (fields map[string]interface{}, err error) {
	err = obj.Unpack(&fields)
	err = errors.Wrap(err, "obj.Unpack failed")
	return
}

// Pack will re-marshal an object from json
func (obj *Object) Pack(v interface{}) (err error) {
	// decode the bytes
	var bytes json.RawMessage
	bytes, err = jsoniter.Marshal(v)
	err = errors.Wrap(err, "jsoniter.Marshal failed")
	if err != nil {
		return
	}

	// and store it as the new source
	obj.Source = &bytes

	return
}
