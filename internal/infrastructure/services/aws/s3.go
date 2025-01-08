package aws

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/aws/smithy-go"
	"github.com/oprimogus/cardapiogo/internal/config"
)

type Region string

const (
	SouthEast1 Region = "sa-east-1"
)

type ClientS3 struct {
	s3 *s3.Client
}

func NewClientS3(s3 *s3.Client) ClientS3 {
	return ClientS3{s3: s3}
}

func (c *ClientS3) CreateBucket(ctx context.Context, bucketName string) error {
	_, err := c.s3.CreateBucket(ctx, &s3.CreateBucketInput{
		Bucket: aws.String(bucketName),
		CreateBucketConfiguration: &types.CreateBucketConfiguration{
			LocationConstraint: types.BucketLocationConstraint(SouthEast1),
		},
	})
	if err != nil {
		slog.ErrorContext(ctx, fmt.Sprintf("Couldn't create bucket %v in Region %v. Here's why: %v\n",
			bucketName, SouthEast1, err))
		return err
	}
	return nil
}

func (c *ClientS3) BucketExists(ctx context.Context, bucketName string) (bool, error) {
	_, err := c.s3.HeadBucket(ctx, &s3.HeadBucketInput{
		Bucket: aws.String(bucketName),
	})
	exists := true
	if err != nil {
		var apiError smithy.APIError
		if errors.As(err, &apiError) {
			var notFound *types.NotFound
			switch {
			case errors.As(apiError, &notFound):
				exists = false
				err = nil
			default:
				slog.ErrorContext(ctx, fmt.Sprintf("Either you don't have access to bucket %v or another error occurred. "+
					"Here's what happened: %v\n", bucketName, err))
			}
		}
	}

	return exists, err
}

func (c *ClientS3) getPublicObjectUrl(bucketName string, region Region, objectKey string) string {
	configInstance := config.GetInstance()
	if configInstance.Api.Environment == string(config.Production) {
		return fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", bucketName, region, objectKey)
	}
	return fmt.Sprintf("https://localhost.localstack.cloud:4566/%s/%s", bucketName, objectKey)
}

func (c *ClientS3) UploadFile(
	ctx context.Context,
	bucketName string,
	objectKey string,
	file []byte) (objectURL string, err error) {

	buffer := bytes.NewBuffer(file)

	uploader := manager.NewUploader(c.s3, func(u *manager.Uploader) {
		u.PartSize = 10 * 1024 * 1024 // 10MiB
	})

	_, err = uploader.Upload(ctx, &s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectKey),
		Body:   buffer, // Alterado para bytes.Buffer
	})
	if err != nil {
		slog.ErrorContext(ctx, fmt.Sprintf("Couldn't upload large object to %v:%v. Here's why: %v\n",
			bucketName, objectKey, err))
		return "", err
	}

	return c.getPublicObjectUrl(bucketName, SouthEast1, objectKey), nil
}

func (c *ClientS3) RemoveFile(ctx context.Context, bucketName, objectKey string) error {
	deleteObjectParams := &s3.DeleteObjectInput{
		Bucket: &bucketName,
		Key:    &objectKey,
	}

	result, err := c.s3.DeleteObject(ctx, deleteObjectParams)
	if err != nil {
		slog.ErrorContext(ctx, fmt.Sprintf("Couldn't delete object %v from S3: %v\n"))
		return err
	}
	fmt.Println(result)

	return nil
}
