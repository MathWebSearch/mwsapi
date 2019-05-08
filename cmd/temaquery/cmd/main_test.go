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

func TestMWSQuery(t *testing.T) {
	integrationtest.MarkSkippable(t)

	tests := []struct {
		name      string
		args      *Args
		assetName string
		wantErr   bool
	}{
		{"count everything", &Args{Expressions: []string{"<mws:qvar>x</mws:qvar>"}, Count: true}, "testdata/count_all.json", false},
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
