package storage

import "context"

type Service interface {
	CreateBucket(ctx context.Context, bucketName string) error
	BucketExists(ctx context.Context, bucketName string) (bool, error)
	UploadFile(ctx context.Context, bucketName string,
		objectKey string, file []byte) (objectURL string, err error)
}
