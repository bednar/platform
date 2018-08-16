package query

import (
	"testing"
	"time"
	"context"
	"github.com/influxdata/platform"
	"github.com/influxdata/platform/mock"
	"github.com/google/go-cmp/cmp"
	"github.com/influxdata/platform/kit/errors"
)

type mockLogger struct {
}

func (*mockLogger) Log(Log) error {
	return nil
}

func newBucketServiceWithOneBucket(bucket platform.Bucket) platform.BucketService {
	bs := mock.NewBucketService()
	bs.FindBucketFn = func(ctx context.Context, bucketFilter platform.BucketFilter) (*platform.Bucket, error) {
		if *bucketFilter.Name == bucket.Name {
			return &bucket, nil
		}

		return nil, errors.New("Unknown bucket")
	}

	return bs
}

func TestAuthorizer_Authorize(t *testing.T) {
	ctx := context.Background()
	now := time.Now().UTC()
	logger := &mockLogger{}

	q := `from(bucket:"my_bucket") |> range(start:-2h) |> yield()`
	spec, err := Compile(ctx, q, now)
	if err != nil {
		t.Errorf("Error compiling query: %v", q)
	}

	// Try to authorize with bucket service with no buckets
	// and no authorization
	auth := platform.Authorization{}
	emptyBucketService := mock.NewBucketService()
	authorizer := NewAuthorizer(emptyBucketService)

	err = authorizer.Authorize(ctx, spec, auth, logger)
	if diagnostic := cmp.Diff("Bucket service returned nil bucket", err.Error()); diagnostic != "" {
		t.Errorf("Authorize message mismatch: -want/+got:\n%v", diagnostic)
	}

	// Try to authorize with a bucket service that knows about one bucket
	// (still no authorization)
	id, _ := platform.IDFromString("DEADBEEF")
	bucketService := newBucketServiceWithOneBucket(platform.Bucket{
		Name: "my_bucket",
		ID: *id,
	})

	authorizer = NewAuthorizer(bucketService)
	err = authorizer.Authorize(ctx, spec, auth, logger)
	if diagnostic := cmp.Diff(`No read permission for bucket: "my_bucket"`, err.Error()); diagnostic != "" {
		t.Errorf("Authorize message mismatch: -want/+got:\n%v", diagnostic)
	}

	// Try to authorize with read permission on bucket
	auth = platform.Authorization{
		Permissions: []platform.Permission{platform.ReadBucketPermission(*id)},
	}

	err = authorizer.Authorize(ctx, spec, auth, logger)
	if err != nil {
		t.Errorf("Expected successful authorization, but got error: \"%v\"", err.Error())
	}
}
