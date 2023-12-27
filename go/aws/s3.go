package aws

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	v4 "github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	AWSConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

type S3 struct {
	client        s3.Client
	presignClient *s3.PresignClient
	uploader      *manager.Uploader
}

func NewS3(region string, config *S3Config) (*S3, error) {
	awsCfg, err := AWSConfig.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatal("Cannot initialize default configuration of AWS")
		return nil, err
	}

	s3Client := s3.NewFromConfig(awsCfg, func(o *s3.Options) {
		o.Region = region
	})

	return &S3{
		client: *s3Client,
		presignClient: s3.NewPresignClient(s3Client, func(po *s3.PresignOptions) {
			po.Expires = time.Duration(24 * int64(time.Hour))
		}),
		uploader: manager.NewUploader(s3Client, func(u *manager.Uploader) {
			u.Concurrency = config.UploadConcurrency
			u.PartSize = int64(config.UploadPartSize)
		}),
	}, nil
}

func (s *S3) ListBucket() ([]types.Bucket, error) {
	output, err := s.client.ListBuckets(context.TODO(), nil)
	if err != nil {
		return nil, err
	}

	return output.Buckets, nil
}

func (s *S3) GetObject(bucketName string, objectKey string, lifetimeSecs int64) (*v4.PresignedHTTPRequest, error) {
	request, err := s.presignClient.PresignGetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectKey),
	}, func(opts *s3.PresignOptions) {
		opts.Expires = time.Duration(lifetimeSecs * int64(time.Second))
	})
	if err != nil {
		log.Printf("Couldn't get a presigned request to %v:%v. Here's why: %v\n", bucketName, objectKey, err)
	}
	return request, err
}

func (s *S3) UploadFile(filePath string, bucketName string, objectKey string) error {
	f, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open file %q, %v", filePath, err)
	}
	defer f.Close()

	obj, err := s.client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectKey),
		Body:   f,
	})
	if err != nil {
		return fmt.Errorf("# Failed to upload file %v", err)
	}

	fmt.Printf("# File uploaded to %v\n", obj)
	return nil
}

func (s *S3) UploadLargeFile(bucket, key string, body io.Reader) (*manager.UploadOutput, error) {
	result, err := s.uploader.Upload(context.Background(), &s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
		Body:   body,
	})
	return result, err
}

func (s *S3) DeleteFile(bucket, key string) error {
	obj, err := s.client.DeleteObject(context.TODO(), &s3.DeleteObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return fmt.Errorf("failed to delete file %v", err)
	}
	fmt.Printf("# s3 file deleted %v\n", obj)
	return nil
}
