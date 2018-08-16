package query

import (
	"github.com/pkg/errors"
	"github.com/influxdata/platform"
	"context"
)

// Authorizer provides a method for ensuring that the buckets accessed by a query spec
// are allowed access by a give Authorization
type Authorizer interface {
	Authorize(ctx context.Context, spec *Spec, auth platform.Authorization, logger Logger) error
}

// NewAuthorizer creates a new Authorizer
func NewAuthorizer(bucketService platform.BucketService) Authorizer {
	return &authorizer{bucketService: bucketService}
}

type authorizer struct {
	bucketService platform.BucketService
}

// Authorize finds all the buckets read and written by the given spec, and ensures that execution is allowed
// given the Authorization.  Returns nil on success, and an error with an appropriate message otherwise.
func (a *authorizer) Authorize(ctx context.Context, spec *Spec, auth platform.Authorization, auditLogger Logger) error {

	readBuckets, writeBuckets, err := spec.BucketsAccessed()

	if err != nil {
		return errors.Wrap(err, "Could not retrieve buckets for query.Spec")
	}

	for _, readBucketFilter := range readBuckets {
		bucket, err := a.bucketService.FindBucket(ctx, readBucketFilter)
		if err != nil {
			return errors.Wrapf(err, "Bucket service error")
		} else if bucket == nil {
			return errors.New("Bucket service returned nil bucket")
		}

		reqPerm := platform.ReadBucketPermission(bucket.ID)
		if ! platform.Allowed(reqPerm, auth.Permissions) {
			return errors.New("No read permission for bucket: \"" + bucket.Name + "\"")
		}
	}

	for _, writeBucketFilter := range writeBuckets {
		bucket, err := a.bucketService.FindBucket(context.Background(), writeBucketFilter)
		if err != nil {
			return errors.Wrapf(err, "Could not find bucket %v", writeBucketFilter)
		}

		reqPerm := platform.WriteBucketPermission(bucket.ID)
		if ! platform.Allowed(reqPerm, auth.Permissions) {
			return errors.New("No write permission for bucket: \"" + bucket.Name + "\"")
		}
	}

	// TODO: log pass/fail decision to audit logger

	return nil
}
