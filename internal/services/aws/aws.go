package aws

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsConf "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/oprimogus/cardapiogo/internal/config"
)

var (
	AwsService *AwsInstance
)

type AwsInstance struct {
	conf aws.Config
	S3   ClientS3
}

func newAwsInstance(ctx context.Context) (awsInstance *AwsInstance, err error) {
	configInstance := config.GetInstance()
	cfg, err := awsConf.LoadDefaultConfig(
		ctx,
		awsConf.WithRegion(configInstance.Aws.Region()),
		awsConf.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(
				configInstance.Aws.AccessKeyID(),
				configInstance.Aws.SecretAccessKey(),
				configInstance.Aws.SessionKey(),
			),
		))
	if err != nil {
		slog.ErrorContext(ctx, fmt.Sprintf("fail on load default configuration from AWS: %s", err))
		return nil, err
	}

	var s3Client *s3.Client
	if configInstance.Api.Environment == string(config.Production) {
		s3Client = s3.NewFromConfig(cfg)
	} else {
		s3Client = s3.NewFromConfig(cfg, func(o *s3.Options) {
			o.BaseEndpoint = aws.String("http://localhost:4566")
			o.UsePathStyle = true
		})
	}

	return &AwsInstance{
		conf: cfg,
		S3:   NewClientS3(s3Client),
	}, nil
}

func GetInstance(ctx context.Context) (awsInstance *AwsInstance, err error) {
	if AwsService == nil {
		return newAwsInstance(ctx)
	}
	return AwsService, nil
}
