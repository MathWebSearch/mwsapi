package result

import (
	"testing"
)

func TestMathFormula_SetURL(t *testing.T) {
	type args struct {
		URL string
	}
	tests := []struct {
		name            string
		args            args
		wantDocumentURL string
		wantLocalID     string
	}{
		{"empty URL", args{URL: ""}, "", ""},
		{"url without '#'", args{URL: "local"}, "", "local"},
		{"url with single '#'", args{URL: "global#local"}, "global", "local"},
		{"url with multiple '#'", args{URL: "global#local#sublocal"}, "global#local", "sublocal"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			formula := &MathFormula{}
			formula.SetURL(tt.args.URL)
			if formula.DocumentURL != tt.wantDocumentURL {
				t.Errorf("MathInfo.SetURL().DocumentURL = %v, want %v", formula.DocumentURL, tt.wantDocumentURL)
			}
			if formula.LocalID != tt.wantLocalID {
				t.Errorf("MathInfo.SetURL().LocalID = %v, want %v", formula.LocalID, tt.wantLocalID)
			}
		})
	}
}

func TestMathFormula_RealMathID(t *testing.T) {
	type fields struct {
		LocalID string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{"empty id", fields{LocalID: ""}, ""},

		{"small legacy id", fields{LocalID: "1"}, "math1"},
		{"large legacy id", fields{LocalID: "123456789"}, "math123456789"},

		{"uuid-like", fields{LocalID: "abcdefghijk"}, "abcdefghijk"},
		{"mixed", fields{LocalID: "123abcdefghijk"}, "123abcdefghijk"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			formula := &MathFormula{
				LocalID: tt.fields.LocalID,
			}
			if got := formula.RealMathID(); got != tt.want {
				t.Errorf("MathInfo.RealMathID() = %v, want %v", got, tt.want)
			}
		})
	}
}
