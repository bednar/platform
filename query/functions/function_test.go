package functions

import (
	"testing"
	"github.com/influxdata/platform/query/values"
	"github.com/influxdata/platform/query"
)

// TestBucketAwareOperationSpec verifies that any functions with parameters
// named "bucket" or "bucketID" implement interface BucketAwareOperationSpec.
func TestBucketAwareOperationSpec(t *testing.T) {
	builtIns, _ := query.BuiltIns()
	for name, v := range builtIns {
		fn, ok := v.(values.Function)
		if !ok {
			continue
		}
		params := fn.Type().Params()
		for k := range params {
			if k == "bucket" || k == "bucketID" {
				opSpec := query.OperationSpecNewFn(query.OperationKind(name))()
				if _, ok := opSpec.(query.BucketAwareOperationSpec); !ok {
					t.Errorf(`Operation "%v" does not implement BucketAwareOperationSpec ` +
						`despite having parameters named "bucket" or "bucketID"`, name)
				}
			}
		}
	}
}
