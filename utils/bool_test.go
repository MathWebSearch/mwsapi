package utils

import (
	"encoding/xml"
	"reflect"
	"testing"
)

func TestBooleanYesNo_MarshalText(t *testing.T) {
	type element struct {
		XMLName xml.Name     `xml:"element"`
		Value   BooleanYesNo `xml:"value,attr"`
	}

	tests := []struct {
		name    string
		element element
		wantXML string
		wantErr bool
	}{
		{"true", element{Value: true}, "<element value=\"yes\"></element>", false},
		{"false", element{Value: false}, "<element value=\"no\"></element>", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotBytes, err := xml.Marshal(tt.element)
			gotXML := string(gotBytes)
			if (err != nil) != tt.wantErr {
				t.Errorf("BooleanYesNo.MarshalText() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotXML, tt.wantXML) {
				t.Errorf("BooleanYesNo.MarshalText() = %v, want %v", gotXML, tt.wantXML)
			}
		})
	}
}

func BenchmarkBooleanYesNo_MarshalText(b *testing.B) {
	byes := BooleanYesNo(true)
	bno := BooleanYesNo(false)

	for n := 0; n < b.N; n++ {
		xml.Marshal(byes)
		xml.Marshal(bno)
	}
}

func TestBooleanYesNo_UnmarshalText(t *testing.T) {
	type element struct {
		XMLName xml.Name     `xml:"element"`
		Value   BooleanYesNo `xml:"value,attr"`
	}

	tests := []struct {
		name        string
		xml         string
		wantElement element
		wantErr     bool
	}{
		{"true", "<element value=\"yes\"></element>", element{XMLName: xml.Name{Space: "", Local: "element"}, Value: true}, false},
		{"false", "<element value=\"no\"></element>", element{XMLName: xml.Name{Space: "", Local: "element"}, Value: false}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var gotElement element
			err := xml.Unmarshal([]byte(tt.xml), &gotElement)
			if (err != nil) != tt.wantErr {
				t.Errorf("BooleanYesNo.UnmarshalText() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotElement, tt.wantElement) {
				t.Errorf("BooleanYesNo.UnmarshalText() = %v, want %v", gotElement, tt.wantElement)
			}
		})
	}
}

func BenchmarkBooleanYesNo_UnmarshalText(b *testing.B) {
	byes := []byte("<element value=\"yes\"></element>")
	bno := []byte("<element value=\"no\"></element>")

	var dest struct {
		XMLName xml.Name     `xml:"element"`
		Value   BooleanYesNo `xml:"value,attr"`
	}
	tdest := &dest

	for n := 0; n < b.N; n++ {
		xml.Unmarshal(byes, tdest)
		xml.Unmarshal(bno, tdest)
	}
}
