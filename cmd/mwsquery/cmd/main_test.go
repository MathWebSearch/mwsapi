package cmd

import (
	"testing"

	"github.com/MathWebSearch/mwsapi/integrationtest"
)

func TestMain(m *testing.M) {
	integrationtest.Main(m, "docker-compose-mwsquery.yml", nil, "http://localhost:8080/")
}

const MWSXMLAnyQuery = "<mws:qvar>x</mws:qvar>"
const MWSXMLPlusQuery = "<m:apply><plus /><mws:qvar>x</mws:qvar><mws:qvar>y</mws:qvar></m:apply>"

func TestMWSQuery(t *testing.T) {
	integrationtest.MarkSkippable(t)

	tests := []struct {
		name      string
		args      *Args
		assetName string
		wantErr   bool
	}{
		{"count all regular elements", &Args{Expressions: []string{MWSXMLAnyQuery}, Count: true}, "testdata/count_all_regular.json", false},
		{"first 10 regular elements", &Args{Expressions: []string{MWSXMLAnyQuery}, From: 0, Size: 10}, "testdata/first_10_regular.json", false},
		{"last 10 regular elements", &Args{Expressions: []string{MWSXMLAnyQuery}, From: 2914, Size: 10}, "testdata/last_10_regular.json", false},

		{"count all mws ids", &Args{Expressions: []string{MWSXMLAnyQuery}, Count: true, MWSIdsOnly: true}, "testdata/count_all_ids.json", false},
		{"first 10 mws ids", &Args{Expressions: []string{MWSXMLAnyQuery}, From: 0, Size: 10, MWSIdsOnly: true}, "testdata/first_10_ids.json", false},
		{"last 10 mws ids", &Args{Expressions: []string{MWSXMLAnyQuery}, From: 735, Size: 10, MWSIdsOnly: true}, "testdata/last_10_ids.json", false},

		{"count all terms of shape x + y", &Args{Expressions: []string{MWSXMLPlusQuery}, Count: true}, "testdata/count_plus_regular.json", false},
		{"first 10 terms of shape x + y", &Args{Expressions: []string{MWSXMLPlusQuery}, From: 0, Size: 10}, "testdata/first_10_plus_sums_regular.json", false},
		{"last 10 terms of shape x + y", &Args{Expressions: []string{MWSXMLPlusQuery}, From: 17, Size: 10}, "testdata/last_10_plus_sums_regular.json", false},

		{"count all ids of shape x + y", &Args{Expressions: []string{MWSXMLPlusQuery}, Count: true, MWSIdsOnly: true}, "testdata/count_plus_ids.json", false},
		{"first 10 ids of shape x + y", &Args{Expressions: []string{MWSXMLPlusQuery}, MWSIdsOnly: true, From: 0, Size: 10}, "testdata/first_10_plus_sums_ids.json", false},
		{"last 10 ids of shape x + y", &Args{Expressions: []string{MWSXMLPlusQuery}, MWSIdsOnly: true, From: 12, Size: 10}, "testdata/last_10_plus_sums_ids.json", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			// fix arguments and normalize
			tt.args.MWSHost = "localhost"
			tt.args.MWSPort = 8080
			tt.args.Normalize = true

			gotRes, err := Main(tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("mwsquery() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			integrationtest.TestJSONAsset(t, "mwsquery()", gotRes, tt.assetName)
		})
	}
}
