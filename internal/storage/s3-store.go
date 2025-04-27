package storage

import (
	"context"
	"mime/multipart"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3StorageAdapter struct {
  Bucket string
  Client *s3.Client
}

func NewS3StorageAdapter(bucket string, client *s3.Client) *S3StorageAdapter {
  return &S3StorageAdapter{
      Bucket: bucket,
      Client: client,
  }
}

func (s *S3StorageAdapter) Upload(ctx context.Context, file multipart.File, filename string, contentType string) (string, error) {
  uploader := manager.NewUploader(s.Client)

  result, err := uploader.Upload(ctx, &s3.PutObjectInput{
      Bucket:      aws.String(s.Bucket),
      Key:         aws.String(filename),
      Body:        file,
      ContentType: aws.String(contentType),
  })
  if err != nil {
      return "", err
  }

  return result.Location, nil
}
