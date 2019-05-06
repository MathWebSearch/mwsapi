package cmd

import (
	"net/http"
	"testing"

	"github.com/MathWebSearch/mwsapi/integrationtest"
)

func TestMain(m *testing.M) {
	integrationtest.Main(m, "elasticquery_elastic", func(client *http.Client) error {
		return integrationtest.LoadElasticSnapshot(client, "http://localhost:9300", "/snapshots/")
	}, "http://localhost:9300/")
}

func TestMWSQuery(t *testing.T) {
	integrationtest.MarkSkippable(t)

	tests := []struct {
		name      string
		args      *Args
		assetName string
		wantErr   bool
	}{
		{"count all elements with 'world' text", &Args{Text: "World", Count: true}, "testdata/count_world.json", false},
		{"run document phase for 'world' text ", &Args{Text: "World", DocumentPhaseOnly: true, From: 0, Size: 8}, "testdata/document_world.json", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// fix arguments and normalize
			tt.args.ElasticHost = "localhost"
			tt.args.ElasticPort = 9300
			tt.args.Normalize = true

			gotRes, err := Main(tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("elasticquery() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			integrationtest.TestJSONAsset(t, "elasticquery()", gotRes, tt.assetName)
		})
	}
}
