package confs3

import (
	"context"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/lifecycle"
)

func (p *S3Client) CreateBucket(bucket string) error {
	ctx := context.Background()
	ok, err := p.BucketExists(bucket)
	if err != nil {
		return err
	}

	if ok {
		return nil
	}
	return p.cli.MakeBucket(ctx, bucket, minio.MakeBucketOptions{
		Region: p.Region,
	})
}

func (p *S3Client) DeleteBucket(bucket string) error {
	ctx := context.Background()
	return p.cli.RemoveBucket(ctx, bucket)
}

func (p *S3Client) ListBuckets() ([]minio.BucketInfo, error) {
	ctx := context.Background()
	return p.cli.ListBuckets(ctx)
}

func (p *S3Client) BucketExists(bucket string) (bool, error) {
	ctx := context.Background()
	return p.cli.BucketExists(ctx, bucket)
}

func (p *S3Client) GetBucketLifecycle() (*lifecycle.Configuration, error) {
	ctx := context.Background()
	return p.cli.GetBucketLifecycle(ctx, p.Bucket)
}

func (p *S3Client) SetBucketLifecycleExpireIn(prefix string, days int) error {
	config, err := NewLifeCycle().ExpiresIn(prefix, days)
	if err != nil {
		return err
	}
	ctx := context.Background()
	return p.cli.SetBucketLifecycle(ctx, p.Bucket, config)
}

func (p *S3Client) SetBucketLifecycleExpireAt(prefix string, date string) error {
	config, err := NewLifeCycle().ExpiresAt(prefix, date)
	if err != nil {
		return err
	}
	ctx := context.Background()
	return p.cli.SetBucketLifecycle(ctx, p.Bucket, config)
}
