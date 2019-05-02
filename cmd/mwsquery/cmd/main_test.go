// +build integration

package cmd

import (
	"os"
	"testing"

	"github.com/MathWebSearch/mwsapi/utils"
)

func TestMain(m *testing.M) {
	if utils.StartIntegrationTest(m, "mwsquery") != 0 {
		os.Exit(1)
	}

	code := m.Run()

	utils.StopIntegrationTest(m, "mwsquery")

	os.Exit(code)
}

func TestMWSQuery(t *testing.T) {
	tests := []struct {
		name      string
		args      *Args
		assetName string
		wantErr   bool
	}{
		{"count all elements", &Args{Expressions: []string{"<mws:qvar>x</mws:qvar>"}, Count: true}, "testdata/count_all.json", false},
		{"first 10 elements", &Args{Expressions: []string{"<mws:qvar>x</mws:qvar>"}, From: 0, Size: 10}, "testdata/first_10.json", false},
		{"last 10 elements", &Args{Expressions: []string{"<mws:qvar>x</mws:qvar>"}, From: 2914, Size: 10}, "testdata/last_10.json", false},
		// TODO: More tests
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			// fix arguments and normalize
			tt.args.MWSHost = "localhost"
			tt.args.MWSPort = 8080
			tt.args.Normalize = true

			gotRes, err := Main(tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("Main() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			utils.CompareJSONAsset(t, "mwsquery()", gotRes, tt.assetName)
		})
	}
}
