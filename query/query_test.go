package query

import "testing"

func TestQuery_Type(t *testing.T) {
	type fields struct {
		Expressions []string
		Text        string
	}
	tests := []struct {
		name   string
		fields fields
		want   Kind
	}{
		{"empty query (1)", fields{Expressions: nil, Text: ""}, EmptyQueryKind},
		{"empty query (2) ", fields{Expressions: []string{}, Text: ""}, EmptyQueryKind},

		{"MWSQuery (1)", fields{Expressions: []string{"<mws:qvar>x</mws:qvar>"}, Text: ""}, MWSQueryKind},
		{"MWSQuery (2)", fields{Expressions: []string{"<mws:qvar>y</mws:qvar>"}, Text: ""}, MWSQueryKind},

		{"MWSQuery (1)", fields{Expressions: []string{}, Text: "hello"}, ElasticQueryKind},
		{"MWSQuery (2)", fields{Expressions: nil, Text: "world"}, ElasticQueryKind},

		{"TemaQuery (1)", fields{Expressions: []string{"<mws:qvar>x</mws:qvar>"}, Text: "hello"}, TemaSearchQueryKind},
		{"TemaQuery (2)", fields{Expressions: []string{"<mws:qvar>y</mws:qvar>"}, Text: "world"}, TemaSearchQueryKind},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q := &Query{
				Expressions: tt.fields.Expressions,
				Text:        tt.fields.Text,
			}
			if got := q.Kind(); got != tt.want {
				t.Errorf("Query.Kind() = %v, want %v", got, tt.want)
			}
		})
	}
}
