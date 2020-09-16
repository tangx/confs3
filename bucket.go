package confs3

import (
	"context"

	"github.com/minio/minio-go/v7"
)

func (p *S3Client) MakeBucket(bucket string, region string) error {
	ctx := context.Background()
	return p.cli.MakeBucket(ctx, bucket, minio.MakeBucketOptions{
		Region: region,
	})
}

func (p *S3Client) RemoveBucket(bucket string) error {
	ctx := context.Background()
	return p.cli.RemoveBucket(ctx, bucket)
}

func (p *S3Client) ListBuckets() ([]minio.BucketInfo, error) {
	ctx := context.Background()
	return p.cli.ListBuckets(ctx)
}
