package query

import (
	"encoding/xml"
	"reflect"
	"testing"

	"github.com/MathWebSearch/mwsapi/utils"
)

func TestMWSQuery_Raw(t *testing.T) {
	type fields struct {
		Expressions []string
		MwsIdsOnly  bool
	}
	type args struct {
		from int64
		size int64
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *RawMWSQuery
	}{
		{"without mws-ids", fields{[]string{"<mws:qvar>x</mws:qvar>"}, false}, args{12, 15}, &RawMWSQuery{
			From: 12,
			Size: 15,

			ReturnTotal:  true,
			OutputFormat: "json",

			Expressions: []*MWSExpression{&MWSExpression{"<mws:qvar>x</mws:qvar>", xml.Name{}}},
		}},

		{"with mws-ids", fields{[]string{"<mws:qvar>x</mws:qvar>"}, true}, args{56, 74}, &RawMWSQuery{
			From: 56,
			Size: 74,

			ReturnTotal:  true,
			OutputFormat: "mws-ids",

			Expressions: []*MWSExpression{&MWSExpression{"<mws:qvar>x</mws:qvar>", xml.Name{}}},
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q := &MWSQuery{
				Expressions: tt.fields.Expressions,
				MwsIdsOnly:  tt.fields.MwsIdsOnly,
			}
			if got := q.Raw(tt.args.from, tt.args.size); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MWSQuery.Raw() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRawMWSQuery_ToXML(t *testing.T) {
	type fields struct {
		From         int64
		Size         int64
		ReturnTotal  utils.BooleanYesNo
		OutputFormat string
		Expressions  []string
	}
	tests := []struct {
		name    string
		fields  fields
		want    string
		wantErr bool
	}{
		{"empty query", fields{}, "<mws:query limitmin=\"0\" answsize=\"0\" totalreq=\"no\" output=\"\" xmlns:mws=\"http://www.mathweb.org/mws/ns\" xmlns:m=\"http://www.w3.org/1998/Math/MathML\"></mws:query>", false},
		{"simple query", fields{5, 10, true, "xml", []string{"<m:limit/>"}}, "<mws:query limitmin=\"5\" answsize=\"10\" totalreq=\"yes\" output=\"xml\" xmlns:mws=\"http://www.mathweb.org/mws/ns\" xmlns:m=\"http://www.w3.org/1998/Math/MathML\"><mws:expr><m:limit/></mws:expr></mws:query>", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := make([]*MWSExpression, len(tt.fields.Expressions))
			for i, s := range tt.fields.Expressions {
				e[i] = &MWSExpression{Term: s}
			}
			q := &RawMWSQuery{
				From:         tt.fields.From,
				Size:         tt.fields.Size,
				ReturnTotal:  tt.fields.ReturnTotal,
				OutputFormat: tt.fields.OutputFormat,
				Expressions:  e,
			}
			got, err := q.ToXML()
			gotString := string(got)
			if (err != nil) != tt.wantErr {
				t.Errorf("Query.ToXML() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotString, tt.want) {
				t.Errorf("Query.ToXML() = %v, want %v", gotString, tt.want)
			}
		})
	}
}
