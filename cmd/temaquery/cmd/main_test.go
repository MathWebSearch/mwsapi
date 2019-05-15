package cmd

import (
	"net/http"
	"testing"

	"github.com/MathWebSearch/mwsapi/integrationtest"
)

func TestMain(m *testing.M) {
	integrationtest.Main(m, "docker-compose-temaquery.yml", func(client *http.Client) error {
		return integrationtest.LoadElasticSnapshot(client, "http://localhost:9400", "/snapshots/")
	}, "http://localhost:8181/", "http://localhost:9400/")
}

const MWSXMLAnyQuery = "<mws:qvar>x</mws:qvar>"
const MWSXMLPlusQuery = "<m:apply><plus /><mws:qvar>x</mws:qvar><mws:qvar>y</mws:qvar></m:apply>"

func TestTemaQuery(t *testing.T) {
	integrationtest.MarkSkippable(t)

	tests := []struct {
		name      string
		args      *Args
		assetName string
		wantErr   bool
	}{
		// count everything
		{"count nothing", &Args{Count: true}, "testdata/count_nothing.json", false},
		{"count everything", &Args{Expressions: []string{MWSXMLAnyQuery}, Count: true}, "testdata/count_all.json", false},

		{"count all terms of the form x + y", &Args{Expressions: []string{MWSXMLPlusQuery}, Count: true}, "testdata/count_plus.json", false},
		{"count all documents with the word 'is'", &Args{Text: "is", Count: true}, "testdata/count_is.json", false},
		{"count all documents with the word 'is' and terms of the form x + y", &Args{Expressions: []string{MWSXMLPlusQuery}, Text: "is", Count: true}, "testdata/count_plus_is.json", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// fix arguments and normalize
			tt.args.MWSHost = "localhost"
			tt.args.MWSPort = 8181
			tt.args.ElasticHost = "localhost"
			tt.args.ElasticPort = 9400
			tt.args.Normalize = true

			gotRes, err := Main(tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("temaquery() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			integrationtest.TestJSONAsset(t, "temaquery()", gotRes, tt.assetName)
		})
	}
}
