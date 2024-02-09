package objectstorage

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3ObjectStorage struct {
	client     *s3.Client
	bucketName string
}

func NewS3ObjectStorage(clientId, clientSecret, bucketName string) (os *S3ObjectStorage, err error) {
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(clientId, clientSecret, "")),
		config.WithRegion("ap-southeast-1"))

	if err != nil {
		return
	}

	client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.UsePathStyle = true
	})

	os = &S3ObjectStorage{
		client:     client,
		bucketName: bucketName,
	}
	return
}

func (os *S3ObjectStorage) ListObject(ctx context.Context) (keys []string, err error) {
	out, err := os.client.ListObjectsV2(ctx, &s3.ListObjectsV2Input{
		Bucket: aws.String(os.bucketName),
	})

	if err != nil {
		return
	}

	for _, o := range out.Contents {
		keys = append(keys, *o.Key)
	}

	return
}

func (os *S3ObjectStorage) GetObject(ctx context.Context, key string) (url string, err error) {
	pClient := s3.NewPresignClient(os.client, func(po *s3.PresignOptions) {
		po.Expires = time.Minute * 15
	})

	out, err := pClient.PresignGetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(os.bucketName),
		Key:    aws.String(key),
	})

	if err != nil {
		return
	}

	url = out.URL

	return
}
