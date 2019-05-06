package integrationtest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
)

// LoadElasticSnapshot loads an elastic snapshot
func LoadElasticSnapshot(client *http.Client, url string, location string) (err error) {
	reponame := "preload"
	snapshotname := "preload"

	if testing.Verbose() {
		fmt.Println("Restoring Elasticsearch snapshot for testing ...")
	}

	err = registerBackupLocation(client, url, reponame, location)
	if err != nil {
		return
	}

	err = loadSnapshot(client, url, reponame, snapshotname)
	if err != nil {
		return
	}

	return refreshElasticSearch(client, url)
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
	d, err := json.Marshal(data)
	if err != nil {
		return
	}

	// make a PUT request to the appropriate url
	uri := fmt.Sprintf("%s/_snapshot/%s", url, reponame)
	req, err := http.NewRequest(http.MethodPut, uri, bytes.NewBuffer(d))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json; charset=utf-8")

	// and send it
	res, err := client.Do(req)
	if err != nil {
		return
	}

	bodyBytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return
	}

	// and return
	if res.StatusCode != 200 {
		return fmt.Errorf("%s returned %s", uri, string(bodyBytes))
	}
	return
}

func loadSnapshot(client *http.Client, url string, reponame string, snapshotname string) (err error) {

	// the body of the request
	data := j{
		"include_global_state": true,
	}
	d, err := json.Marshal(data)
	if err != nil {
		return
	}

	// the request itself
	uri := fmt.Sprintf("%s/_snapshot/%s/%s/_restore", url, reponame, snapshotname)
	req, err := http.NewRequest(http.MethodPost, uri, bytes.NewBuffer(d))
	req.Header.Set("Content-Type", "application/json; charset=utf-8")

	if err != nil {
		return
	}

	// and send it
	res, err := client.Do(req)
	if err != nil {
		return
	}

	return expectHTTP200(res, uri)
}

func refreshElasticSearch(client *http.Client, url string) (err error) {
	uri := fmt.Sprintf("%s/_refresh", url)
	res, err := client.Get(uri)
	if err != nil {
		return
	}

	return expectHTTP200(res, uri)
}

func expectHTTP200(res *http.Response, uri string) (err error) {
	bodyBytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return
	}

	if res.StatusCode != 200 {
		return fmt.Errorf("%s returned %s", uri, string(bodyBytes))
	}

	return
}
