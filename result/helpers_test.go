package result

import (
	"testing"
)

func TestMathInfo_ID(t *testing.T) {
	type fields struct {
		Source string
		URL    string
		XPath  string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{"empty URL", fields{URL: ""}, ""},
		{"url without '#'", fields{URL: "local"}, "local"},
		{"url with single '#'", fields{URL: "global#local"}, "local"},
		{"url with multiple '#'", fields{URL: "global#local#sublocal"}, "sublocal"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			info := &MathInfo{
				Source: tt.fields.Source,
				URL:    tt.fields.URL,
				XPath:  tt.fields.XPath,
			}
			if got := info.ID(); got != tt.want {
				t.Errorf("MathInfo.ID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMathInfo_RealMathID(t *testing.T) {
	type fields struct {
		Source string
		URL    string
		XPath  string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{"empty id", fields{URL: ""}, ""},

		{"small legacy id", fields{URL: "1"}, "math1"},
		{"large legacy id", fields{URL: "123456789"}, "math123456789"},

		{"uuid-like", fields{URL: "abcdefghijk"}, "abcdefghijk"},
		{"mixed", fields{URL: "123abcdefghijk"}, "123abcdefghijk"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			info := &MathInfo{
				Source: tt.fields.Source,
				URL:    tt.fields.URL,
				XPath:  tt.fields.XPath,
			}
			if got := info.RealMathID(); got != tt.want {
				t.Errorf("MathInfo.RealMathID() = %v, want %v", got, tt.want)
			}
		})
	}
}
