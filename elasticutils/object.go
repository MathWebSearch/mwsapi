package elasticutils

import (
	"context"
	"encoding/json"
	"errors"

	"gopkg.in/olivere/elastic.v6"
)

// Object represents an object within elasticsearch
type Object struct {

	// the client that this object uses
	client *elastic.Client

	index string // index the object resides in
	tp    string // type the object resides in
	id    string // id of this object

	Fields map[string]interface{} // the fields of this object

	Hit *HitInfo // information about this object as a result (if available)
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
		return errors.New("Already indexed")
	}

	// perform the Indexing operation
	ctx := context.Background()
	result, err := obj.client.Index().Index(obj.index).Type(obj.tp).BodyJson(obj.Fields).Do(ctx)
	if err == nil && result.Shards.Successful <= 0 {
		err = errors.New("Index() reported 0 successful shards")
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
		return errors.New("Not indexed")
	}

	ctx := context.Background()

	// fetch from the db and return unless an error occured
	result, err := obj.client.Get().Index(obj.index).Type(obj.tp).Id(obj.id).Do(ctx)
	if err == nil && !result.Found {
		err = errors.New("Reload() did not find object")
	}

	if err != nil {
		err = obj.setSource(result.Source)
	}

	return
}

func (obj *Object) updateFromHit(source *elastic.SearchHit) error {
	obj.Hit = &HitInfo{
		Score:     source.Score,
		Highlight: &source.Highlight,
	}
	return obj.setSource(source.Source)
}

func (obj *Object) setSource(source *json.RawMessage) error {
	return json.Unmarshal(*source, &obj.Fields)
}

// Save saves this object into the database
func (obj *Object) Save() (err error) {
	if !obj.IsIndexed() {
		return errors.New("Not indexed")
	}

	// replace the entire item in the database
	ctx := context.Background()
	res, err := obj.client.Update().Index(obj.index).Type(obj.tp).Id(obj.id).Doc(obj.Fields).Do(ctx)

	if err == nil && (res.Result != "noop" && res.Shards.Successful <= 0) {
		err = errors.New("Save() reported non-noop result with 0 successful shards ")
	}

	return
}

// Delete deletes this object
func (obj *Object) Delete() (err error) {
	if !obj.IsIndexed() {
		return errors.New("Not indexed")
	}

	// just clears the object
	ctx := context.Background()
	res, err := obj.client.Delete().Index(obj.index).Type(obj.tp).Id(obj.id).Do(ctx)

	if err == nil && res.Result != "deleted" {
		err = errors.New("Delete() did not report deleted result ")
	}

	// id no longer valid
	if err == nil {
		obj.id = ""
	}

	return
}

// Unpack will unpack this object as json
func (obj *Object) Unpack(v interface{}) (err error) {
	return Repack(obj.Fields, v)
}

// Pack will re-marshal an object from json
func (obj *Object) Pack(v interface{}) (err error) {
	return Repack(v, &obj.Fields)
}

// Repack repacks an interface by going via json
// to should be a pointer
func Repack(from interface{}, to interface{}) error {
	// marshal the interface
	bytes, err := json.Marshal(from)
	if err != nil {
		return err
	}

	// and unmarshal again
	return json.Unmarshal(bytes, to)
}
