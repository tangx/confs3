package confs3

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"path/filepath"

	"github.com/cloudflare/cfssl/log"
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
		_, err := p.StatObject(object)
		if err != nil {
			return p.cli.PresignedPutObject(ctx, p.Bucket, object, p.PresignExpiresIn)
		}
		return nil, fmt.Errorf("%s object already exists", object)
	}

	return p.cli.PresignedPutObject(ctx, p.Bucket, object, p.PresignExpiresIn)
}

func (p *S3Client) StatObject(object string) (minio.ObjectInfo, error) {
	ctx := context.Background()
	return p.cli.StatObject(ctx, p.Bucket, object, minio.StatObjectOptions{})
}

func (p *S3Client) DeleteObject(object string) error {
	ctx := context.Background()
	return p.cli.RemoveObject(ctx, p.Bucket, object, minio.RemoveObjectOptions{})
}

func (p *S3Client) UploadFile(dest string, src string, force bool) (minio.UploadInfo, error) {
	info, err := p.StatObject(dest)
	if err != nil {
		return minio.UploadInfo{}, err
	}

	if info.ETag != "" && !force {
		fmt.Printf("%s Exist, skip\n", dest)
		return minio.UploadInfo{}, nil
	}

	ctx := context.Background()
	n := uint(5)

	fmt.Printf("copy %s -> %s\n", src, filepath.Join(p.Bucket, dest))
	return p.cli.FPutObject(ctx, p.Bucket, dest, src, minio.PutObjectOptions{
		NumThreads: n,
	})
}

func (p *S3Client) UploadFolder(dest string, src string, recursive bool, force bool) error {
	finfo, err := os.Stat(src)
	if err != nil {
		return err
	}
	if !finfo.IsDir() {
		return fmt.Errorf("not a folder")
	}

	fs, err := ioutil.ReadDir(src)
	if err != nil {
		return nil
	}

	for _, f := range fs {
		subsrc := filepath.Join(src, f.Name())
		subdest := filepath.Join(dest, f.Name())

		if !f.Mode().IsDir() {
			_, err = p.UploadFile(subdest, subsrc, force)
			if err != nil {
				log.Warning(err)
			}
		}

		if f.IsDir() && recursive {
			_ = p.UploadFolder(subdest, subsrc, recursive, force)
		}
	}

	return nil
}
