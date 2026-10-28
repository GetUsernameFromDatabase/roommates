// Package storage keeps file contents in an S3 compatible object store
package storage

import (
	"context"
	"errors"
	"io"
	g "roommates/internal/config"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

type Store interface {
	// Put stores size bytes of body under key, replacing what was there
	Put(ctx context.Context, key string, body io.ReadSeeker, size int64, mime string) error
	// Get returns the object, the caller closes it
	Get(ctx context.Context, key string) (io.ReadCloser, error)
	// Delete succeeds when the object does not exist
	Delete(ctx context.Context, key string) error
}

// S3Store is a [Store] on one bucket
type S3Store struct {
	client *s3.Client
	bucket string
}

// New returns an [S3Store] from env, see [g.EVK_S3Endpoint], and creates the bucket if it is missing
//
// panics if the env is missing or the store cannot be reached
func New(ctx context.Context) *S3Store {
	//exhaustruct:ignore
	client := s3.New(s3.Options{
		BaseEndpoint: aws.String(g.EVK_S3Endpoint.MustGet()),
		Region:       g.EVK_S3Region.Get("us-east-1"),
		Credentials: credentials.NewStaticCredentialsProvider(
			g.EVK_S3AccessKey.MustGet(), g.EVK_S3SecretKey.MustGet(), ""),
		UsePathStyle: true,
		// the default CRC checksums are AWS only, other stores may reject them
		RequestChecksumCalculation: aws.RequestChecksumCalculationWhenRequired,
		ResponseChecksumValidation: aws.ResponseChecksumValidationWhenRequired,
	})
	s := &S3Store{client: client, bucket: g.EVK_S3Bucket.MustGet()}
	if err := s.ensureBucket(ctx); err != nil {
		panic("object storage: " + err.Error())
	}
	return s
}

func (s *S3Store) ensureBucket(ctx context.Context) error {
	//exhaustruct:ignore
	_, err := s.client.HeadBucket(ctx, &s3.HeadBucketInput{Bucket: &s.bucket})
	var notFound *types.NotFound
	if !errors.As(err, &notFound) {
		return err
	}
	//exhaustruct:ignore
	_, err = s.client.CreateBucket(ctx, &s3.CreateBucketInput{Bucket: &s.bucket})
	return err
}

func (s *S3Store) Put(ctx context.Context, key string, body io.ReadSeeker, size int64, mime string) error {
	//exhaustruct:ignore
	_, err := s.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:        &s.bucket,
		Key:           &key,
		Body:          body,
		ContentLength: &size,
		ContentType:   &mime,
	})
	return err
}

func (s *S3Store) Get(ctx context.Context, key string) (io.ReadCloser, error) {
	//exhaustruct:ignore
	out, err := s.client.GetObject(ctx, &s3.GetObjectInput{Bucket: &s.bucket, Key: &key})
	if err != nil {
		return nil, err
	}
	return out.Body, nil
}

func (s *S3Store) Delete(ctx context.Context, key string) error {
	//exhaustruct:ignore
	_, err := s.client.DeleteObject(ctx, &s3.DeleteObjectInput{Bucket: &s.bucket, Key: &key})
	return err
}
