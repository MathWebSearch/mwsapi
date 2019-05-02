// +build integration

package cmd

import (
	"os"
	"reflect"
	"testing"

	"github.com/MathWebSearch/mwsapi/utils"
)

func TestMain(m *testing.M) {
	// exit code
	var code int
	defer os.Exit(code)

	// start and stop the integration tests
	code = utils.StartIntegrationTest(m, "mwsquery")
	if code != 0 {
		return
	}
	defer utils.StopIntegrationTest(m, "mwsquery")

	code = m.Run()
	return
}

func TestMWSQuery(t *testing.T) {
	tests := []struct {
		name    string
		args    *Args
		wantRes interface{}
		wantErr bool
	}{
		{"count query", &Args{Expressions: []string{"<mws:qvar>x</mws:qvar>"}, Count: true}, int64(2924), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup host and port
			tt.args.MWSHost = "localhost"
			tt.args.MWSPort = 1000

			gotRes, err := Main(tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("Main() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Main() = %v, want %v", gotRes, tt.wantRes)
			}
		})
	}
}
