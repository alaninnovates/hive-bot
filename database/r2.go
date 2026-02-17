package database

import (
	"context"
	"fmt"
	"io"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type R2 struct {
	client     *s3.Client
	bucketName string
}

func NewR2() *R2 {
	return &R2{bucketName: "hive-images"}
}

func (r2 *R2) Connect(accountId string, accessKeyId string, accessKeySecret string) error {
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessKeyId, accessKeySecret, "")),
		config.WithRegion("auto"),
	)
	if err != nil {
		return err
	}

	r2.client = s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.BaseEndpoint = aws.String(fmt.Sprintf("https://%s.r2.cloudflarestorage.com", accountId))
		o.RequestChecksumCalculation = aws.RequestChecksumCalculationWhenRequired
		o.ResponseChecksumValidation = aws.ResponseChecksumValidationWhenRequired
	})
	return nil
}

func (r2 *R2) UploadImage(key string, data *io.PipeReader) (string, error) {
	uploader := manager.NewUploader(r2.client)

	_, err := uploader.Upload(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(r2.bucketName),
		Key:    aws.String(key),
		Body:   data,
	})
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("https://hive-images.meta-bee.com/%s", key), nil
}
