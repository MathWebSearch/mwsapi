package utils

import (
	"io/ioutil"
	"reflect"
	"testing"
)

func TestMathML(t *testing.T) {
	tests := []struct {
		name string

		assetIn  string
		assetOut string

		wantParseErr bool

		navigate        string
		wantNavigateErr bool
	}{
		{"variable-only math parse", "testdata/mathml/in1.xml", "testdata/mathml/out1.xml", false, "", false},
		{"more complex math parse", "testdata/mathml/in2.xml", "testdata/mathml/out2.xml", false, "", false},
		{"more-complex math parse and navigate", "testdata/mathml/in3.xml", "testdata/mathml/out3.xml", false, "/*[2]/*[5]", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// read input
			source, err := ioutil.ReadFile(tt.assetIn)
			if err != nil {
				t.Error(err.Error())
				return
			}

			// read output
			wantXML, err := ioutil.ReadFile(tt.assetOut)
			if err != nil {
				t.Error(err.Error())
				return
			}

			// read source
			gotMath, err := ParseMathML(string(source))
			if (err != nil) != tt.wantParseErr {
				t.Errorf("ParseMathML() error = %v, wantErr %v", err, tt.wantParseErr)
				return
			}

			// navigate
			err = gotMath.NavigateAnnotation(tt.navigate)
			if (err != nil) != tt.wantNavigateErr {
				t.Errorf("ParseMathML().NavigateSemantic() error = %v, wantErr %v", err, tt.wantNavigateErr)
				return
			}

			// write it again
			gotXML := gotMath.OutputXML()
			if !reflect.DeepEqual(gotXML, string(wantXML)) {
				t.Errorf("ParseMathML().OutputXML() = %v, want %v", gotXML, string(wantXML))
			}
		})
	}
}
