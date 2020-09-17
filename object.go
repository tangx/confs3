package confs3

import (
	"context"
	"fmt"
	"net/url"
	"path/filepath"

	"github.com/minio/minio-go/v7"
)

func (p *S3Client) PreSignedGetURL(object string) (*url.URL, error) {
	ctx := context.Background()
	req := make(url.Values)
	fname := filepath.Base(object)
	req.Set(`response-content-disposition`, fmt.Sprintf(`attachment; filename="%s"`, fname))
	return p.cli.PresignedGetObject(ctx, p.Bucket, object, p.PresignExpiresIn, req)
}

// PreSignedPutURL return url.Value, if not force will return an error when object exists
func (p *S3Client) PreSignedPutURL(object string, force bool) (*url.URL, error) {
	ctx := context.Background()
	if !force {
		_, err := p.StateObject(object)
		if err != nil {
			return p.cli.PresignedPutObject(ctx, p.Bucket, object, p.PresignExpiresIn)
		}
		return nil, fmt.Errorf("%s object already exists", object)
	}

	return p.cli.PresignedPutObject(ctx, p.Bucket, object, p.PresignExpiresIn)
}

func (p *S3Client) StateObject(object string) (minio.ObjectInfo, error) {
	ctx := context.Background()
	return p.cli.StatObject(ctx, p.Bucket, object, minio.StatObjectOptions{})
}

func (p *S3Client) DeleteObject(object string) error {
	ctx := context.Background()
	return p.cli.RemoveObject(ctx, p.Bucket, object, minio.RemoveObjectOptions{})
}

func (p *S3Client) UploadFile(dest string, src string) (minio.UploadInfo, error) {
	ctx := context.Background()
	n := uint(5)
	return p.cli.FPutObject(ctx, p.Bucket, dest, src, minio.PutObjectOptions{
		NumThreads: n,
	})
}

func (p *S3Client) UploadFolder(dest string, src string, recursive bool) {

}
