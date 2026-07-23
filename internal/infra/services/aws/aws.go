package aws

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsConf "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/oprimogus/flyfood-api/internal/config"
)

var (
	Service *Instance
)

type Instance struct {
	conf aws.Config
	S3   ClientS3
}

func newAwsInstance(ctx context.Context) (awsInstance *Instance, err error) {
	configInstance := config.Get()
	cfg, err := awsConf.LoadDefaultConfig(
		ctx,
		awsConf.WithRegion(configInstance.AWS.Region),
		awsConf.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(
				configInstance.AWS.AccessKeyID,
				configInstance.AWS.SecretAccessKey,
				configInstance.AWS.SessionKey,
			),
		))
	if err != nil {
		slog.ErrorContext(ctx, fmt.Sprintf("fail on load default configuration from AWS: %s", err))
		return nil, err
	}

	var s3Client *s3.Client
	if configInstance.API.Environment == string(config.Production) {
		s3Client = s3.NewFromConfig(cfg)
	} else {
		s3Client = s3.NewFromConfig(cfg, func(o *s3.Options) {
			o.BaseEndpoint = aws.String("http://localhost:4566")
			o.UsePathStyle = true
		})
	}

	return &Instance{
		conf: cfg,
		S3:   NewClientS3(s3Client),
	}, nil
}

func GetInstance(ctx context.Context) (awsInstance *Instance, err error) {
	if Service == nil {
		return newAwsInstance(ctx)
	}
	return Service, nil
}
