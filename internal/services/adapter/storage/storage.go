package storage

import "context"

type Bucket string

const (
	BucketProfileImage Bucket = "cardapiogo-profile-images"
	BucketHeaderImage  Bucket = "cardapiogo-header-images"
)

type Service interface {
	CreateBucket(ctx context.Context, bucketName Bucket) error
	BucketExists(ctx context.Context, bucketName Bucket) (bool, error)
	UploadFile(ctx context.Context, bucketName Bucket,
		objectKey string, file []byte) (objectURL string, err error)
}
