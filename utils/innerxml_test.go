package utils

import (
	"encoding/xml"
	"reflect"
	"testing"
)

func TestInnerXML_ToXML(t *testing.T) {
	type args struct {
		XMLName xml.Name `xml:"args"`
		Text    InnerXML `xml:"element"`
	}
	tests := []struct {
		name    string
		args    args
		wantXML string
		wantErr bool
	}{
		{"empty text", args{Text: ""}, "<args><element></element></args>", false},
		{"normal text", args{Text: "Hello World"}, "<args><element>Hello World</element></args>", false},
		{"xml-like text", args{Text: "<hello></hello>"}, "<args><element><hello></hello></element></args>", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := xml.Marshal(&tt.args)
			gotString := string(got)
			if (err != nil) != tt.wantErr {
				t.Errorf("InnerXML.ToXML() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotString, tt.wantXML) {
				t.Errorf("InnerXML.ToXML() = %v, want %v", gotString, tt.wantXML)
			}
		})
	}
}

func TestInnerXML_FromXML(t *testing.T) {
	type args struct {
		Text InnerXML `xml:"element"`
	}
	tests := []struct {
		name     string
		xml      string
		wantArgs args
		wantErr  bool
	}{
		{"empty text", "<args><element></element></args>", args{Text: ""}, false},
		{"normal text", "<args><element>Hello World</element></args>", args{Text: "Hello World"}, false},
		{"xml-like text", "<args><element><hello></hello></element></args>", args{Text: "<hello></hello>"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got args
			err := xml.Unmarshal([]byte(tt.xml), &got)
			if (err != nil) != tt.wantErr {
				t.Errorf("InnerXML.FromXML() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.wantArgs) {
				t.Errorf("InnerXML.FromXML() = %v, want %v", got, tt.wantArgs)
			}
		})
	}
}
