package platform

import "testing"

func TestInternalBucketID(t *testing.T) {
	o := &Organization{ID: []byte("abc")}
	expected := "bar"
	bucketName := InternalBucketName(o)
	if bucketName != expected {
		t.Errorf("Internal bucket name incorrect, got: %s, expected: %s.", bucketName, expected)
	}
}
