package integrationtest

import (
	"bytes"
	"github.com/json-iterator/go"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/pkg/errors"
)

// LoadElasticSnapshot loads an elastic snapshot
func LoadElasticSnapshot(client *http.Client, url string, location string) (err error) {
	reponame := "preload"
	snapshotname := "preload"

	if testing.Verbose() {
		fmt.Println("Restoring Elasticsearch snapshot for testing ...")
	}

	err = registerBackupLocation(client, url, reponame, location)
	err = errors.Wrap(err, "registerBackupLocation failed")
	if err != nil {
		return
	}

	err = loadSnapshot(client, url, reponame, snapshotname)
	err = errors.Wrap(err, "loadSnapshot failed")
	if err != nil {
		return
	}

	err = refreshElasticSearch(client, url)
	err = errors.Wrap(err, "refreshElasticSearch failed")
	return
}

type j map[string]interface{}

func registerBackupLocation(client *http.Client, url string, reponame string, location string) (err error) {

	// the body of the request
	data := j{
		"type": "fs",
		"settings": j{
			"location": location,
			"readonly": true,
		},
	}
	d, err := jsoniter.Marshal(data)
	err = errors.Wrap(err, "jsoniter.Marshal failed")
	if err != nil {
		return
	}

	// make a PUT request to the appropriate url
	uri := fmt.Sprintf("%s/_snapshot/%s", url, reponame)
	req, err := http.NewRequest(http.MethodPut, uri, bytes.NewBuffer(d))
	err = errors.Wrap(err, "http.NewRequest failed")
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json; charset=utf-8")

	// and send it
	res, err := client.Do(req)
	err = errors.Wrap(err, "client.Do failed")
	if err != nil {
		return
	}

	bodyBytes, err := ioutil.ReadAll(res.Body)
	err = errors.Wrap(err, "ioutil.ReadAll failed")
	if err != nil {
		return
	}

	// and return
	if res.StatusCode != 200 {
		return errors.Errorf("%s returned %s", uri, string(bodyBytes))
	}
	return
}

func loadSnapshot(client *http.Client, url string, reponame string, snapshotname string) (err error) {

	// the body of the request
	data := j{
		"include_global_state": true,
	}
	d, err := jsoniter.Marshal(data)
	err = errors.Wrap(err, "jsoniter.Marshal failed")
	if err != nil {
		return
	}

	// the request itself
	uri := fmt.Sprintf("%s/_snapshot/%s/%s/_restore", url, reponame, snapshotname)
	req, err := http.NewRequest(http.MethodPost, uri, bytes.NewBuffer(d))
	err = errors.Wrap(err, "http.NewRequest failed")
	if err != nil {
		return
	}

	// set the header
	req.Header.Set("Content-Type", "application/json; charset=utf-8")

	// and send it
	res, err := client.Do(req)
	err = errors.Wrap(err, "client.Do failed")
	if err != nil {
		return
	}

	err = expectHTTP200(res, uri)
	err = errors.Wrap(err, "expectHTTP200 failed")
	return
}

func refreshElasticSearch(client *http.Client, url string) (err error) {
	uri := fmt.Sprintf("%s/_refresh", url)
	res, err := client.Get(uri)
	err = errors.Wrap(err, "client.Get failed")
	if err != nil {
		return
	}

	err = expectHTTP200(res, uri)
	err = errors.Wrap(err, "expectHTTP200 failed")
	return
}

func expectHTTP200(res *http.Response, uri string) (err error) {
	bodyBytes, err := ioutil.ReadAll(res.Body)
	err = errors.Wrap(err, "ioutil.ReadAll failed")
	if err != nil {
		return
	}

	if res.StatusCode != 200 {
		return errors.Errorf("%s returned %s", uri, string(bodyBytes))
	}

	return
}
