package cmd

import (
	"net/http"
	"testing"

	"github.com/MathWebSearch/mwsapi/integrationtest"
)

func TestMain(m *testing.M) {
	integrationtest.Main(m, "docker-compose-elasticquery.yml", func(client *http.Client) error {
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
		// text only
		{"count all elements with 'world' text", &Args{Text: "World", Count: true}, "testdata/count_world.json", false},
		{"run document phase for 'world' text ", &Args{Text: "World", DocumentPhaseOnly: true, From: 0, Size: 8}, "testdata/document_world.json", false},
		{"run document + highlight phase for 'world' text ", &Args{Text: "World", From: 0, Size: 8}, "testdata/highlight_world.json", false},

		// ids only
		{"count all elements ids 15/53/173", &Args{IDs: []int64{15, 53, 173}, Count: true}, "testdata/count_ids.json", false},
		{"run document phase for all elements with ids 15/53/173 ", &Args{IDs: []int64{15, 53, 173}, DocumentPhaseOnly: true, From: 0, Size: 10}, "testdata/document_ids.json", false},
		{"run document + highlight phase for all elements with ids 15/53/173", &Args{IDs: []int64{15, 53, 173}, From: 0, Size: 10}, "testdata/highlight_ids.json", false},

		// text + ids
		{"count all in text + ids query", &Args{Text: "World", IDs: []int64{15, 53, 173}, Count: true}, "testdata/count_all.json", false},
		{"run document phase for text + ids query", &Args{Text: "World", IDs: []int64{15, 53, 173}, DocumentPhaseOnly: true, From: 0, Size: 8}, "testdata/document_all.json", false},
		{"run document + highlight phase for text + ids query", &Args{Text: "World", IDs: []int64{15, 53, 173}, From: 0, Size: 8}, "testdata/highlight_all.json", false},
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
