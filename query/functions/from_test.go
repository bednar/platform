package functions_test

import (
	"testing"
	"time"

	"github.com/influxdata/platform/query"
	"github.com/influxdata/platform/query/execute"
	"github.com/influxdata/platform/query/functions"
	"github.com/influxdata/platform/query/querytest"
	platformtesting "github.com/influxdata/platform/testing"
)

func TestFrom_NewQuery(t *testing.T) {
	tests := []querytest.NewQueryTestCase{
		{
			Name:    "from no args",
			Raw:     `from()`,
			WantErr: true,
		},
		{
			Name:    "from conflicting args",
			Raw:     `from(db:"d", bucket:"b")`,
			WantErr: true,
		},
		{
			Name:    "from repeat arg",
			Raw:     `from(db:"telegraf", db:"oops")`,
			WantErr: true,
		},
		{
			Name:    "from",
			Raw:     `from(db:"telegraf", chicken:"what is this?")`,
			WantErr: true,
		},
		{
			Name:    "from bucket invalid ID",
			Raw:     `from(bucketID:"invalid")`,
			WantErr: true,
		},
		{
			Name: "from bucket ID",
			Raw:  `from(bucketID:"aaaaaaaaaaaaaaaa")`,
			Want: &query.Spec{
				Operations: []*query.Operation{
					{
						ID: "from0",
						Spec: &functions.FromOpSpec{
							BucketID: platformtesting.MustIDFromString("aaaaaaaaaaaaaaaa"),
						},
					},
				},
			},
		},
		{
			Name: "from with database",
			Raw:  `from(db:"mydb") |> range(start:-4h, stop:-2h) |> sum()`,
			Want: &query.Spec{
				Operations: []*query.Operation{
					{
						ID: "from0",
						Spec: &functions.FromOpSpec{
							Database: "mydb",
						},
					},
					{
						ID: "range1",
						Spec: &functions.RangeOpSpec{
							Start: query.Time{
								Relative:   -4 * time.Hour,
								IsRelative: true,
							},
							Stop: query.Time{
								Relative:   -2 * time.Hour,
								IsRelative: true,
							},
							TimeCol:  "_time",
							StartCol: "_start",
							StopCol:  "_stop",
						},
					},
					{
						ID: "sum2",
						Spec: &functions.SumOpSpec{
							AggregateConfig: execute.DefaultAggregateConfig,
						},
					},
				},
				Edges: []query.Edge{
					{Parent: "from0", Child: "range1"},
					{Parent: "range1", Child: "sum2"},
				},
			},
		},
	}
	for _, tc := range tests {
		tc := tc
		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()
			querytest.NewQueryTestHelper(t, tc)
		})
	}
}

func TestFromOperation_Marshaling(t *testing.T) {
	data := []byte(`{"id":"from","kind":"from","spec":{"db":"mydb"}}`)
	op := &query.Operation{
		ID: "from",
		Spec: &functions.FromOpSpec{
			Database: "mydb",
		},
	}
	querytest.OperationMarshalingTestHelper(t, data, op)
}
