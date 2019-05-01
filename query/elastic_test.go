package query

import (
	"reflect"
	"testing"

	"gopkg.in/olivere/elastic.v6"
)

func TestElasticQuery_RawDocumentQuery(t *testing.T) {
	type fields struct {
		MathWebSearchIDs []int64
		Text             string
	}
	tests := []struct {
		name    string
		fields  fields
		want    RawElasticQuery
		wantErr bool
	}{
		{"empty query", fields{MathWebSearchIDs: nil, Text: ""}, nil, true},

		{"text-only query", fields{MathWebSearchIDs: nil, Text: "Hello world"}, mustClauseResult(elastic.NewMatchQuery("text", "Hello world").MinimumShouldMatch("2").Operator("or")), false},
		{"id-only query", fields{MathWebSearchIDs: []int64{1, 2, 3}, Text: ""}, mustClauseResult(elastic.NewTermsQuery("mws_ids", int64(1), int64(2), int64(3))), false},
		{"full query", fields{MathWebSearchIDs: []int64{1, 2, 3}, Text: "Hello world"}, mustClauseResult(elastic.NewMatchQuery("text", "Hello world").MinimumShouldMatch("2").Operator("or"), elastic.NewTermsQuery("mws_ids", int64(1), int64(2), int64(3))), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q := &ElasticQuery{
				MathWebSearchIDs: tt.fields.MathWebSearchIDs,
				Text:             tt.fields.Text,
			}
			got, err := q.RawDocumentQuery()
			if (err != nil) != tt.wantErr {
				t.Errorf("ElasticQuery.RawDocumentQuery() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ElasticQuery.RawDocumentQuery() = %v, want %v", got, tt.want)
			}
		})
	}
}

// mustClauseResult builds a new MustClause query for testing
func mustClauseResult(qs ...elastic.Query) elastic.Query {
	query := elastic.NewBoolQuery()
	for _, q := range qs {
		query = query.Must(q)
	}
	return query
}
