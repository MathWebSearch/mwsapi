package cmd

import (
	"testing"

	"github.com/MathWebSearch/mwsapi/integrationtest"
)

func TestMain(m *testing.M) {
	integrationtest.Main(m, "docker-compose-elasticsync.yml", nil, "http://localhost:9200/")
}

func TestMWSQuery(t *testing.T) {
	integrationtest.MarkSkippable(t)

	tests := []struct {
		name      string
		args      *Args
		assetName string
		wantErr   bool
	}{
		{"initial sync", &Args{IndexDir: "testdata/a"}, "testdata/001_initial.json", false},
		{"repatead sync", &Args{IndexDir: "testdata/a"}, "testdata/002_second.json", false},

		{"switch data folder", &Args{IndexDir: "testdata/b"}, "testdata/003_upgrade.json", false},
		{"repatead sync", &Args{IndexDir: "testdata/b"}, "testdata/004_unchanged.json", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// fix arguments and normalize
			tt.args.ElasticHost = "localhost"
			tt.args.ElasticPort = 9200
			tt.args.Normalize = true
			tt.args.Quiet = true

			gotRes, err := Main(tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("mwsquery() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			integrationtest.TestJSONAsset(t, "mwsquery()", gotRes, tt.assetName)
		})
	}
}
