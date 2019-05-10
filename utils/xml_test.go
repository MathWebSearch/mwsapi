package utils

import (
	"reflect"
	"testing"
)

func TestResolveXPath(t *testing.T) {
	type args struct {
		fragment string
		xpath    string
	}
	tests := []struct {
		name        string
		args        args
		wantResults []string
		wantErr     bool
	}{
		// clean xml
		{"Finding one element by name", args{fragment: "<root>Hello world</root>", xpath: "/root"}, []string{"<root>Hello world</root>"}, false},
		{"Finding multiple elements by name", args{fragment: "<person><name>Mathew</name><name>Mathilda</name></person>", xpath: "//name/text()"}, []string{"Mathew", "Mathilda"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotResults, err := ResolveXPath(tt.args.fragment, tt.args.xpath)
			if (err != nil) != tt.wantErr {
				t.Errorf("ResolveXPath() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotResults, tt.wantResults) {
				t.Errorf("ResolveXPath() = %v, want %v", gotResults, tt.wantResults)
			}
		})
	}
}
