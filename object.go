package confs3

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"path/filepath"

	"github.com/minio/minio-go/v7"
	"github.com/sirupsen/logrus"
)

func (p *S3Client) PreSignedGetURL(object string) (*url.URL, error) {
	ctx := context.Background()
	object = p.fullObject(object)

	req := make(url.Values)
	fname := filepath.Base(object)
	req.Set(`response-content-disposition`, fmt.Sprintf(`attachment; filename="%s"`, fname))
	return p.cli.PresignedGetObject(ctx, p.Bucket, object, p.PresignExpiresIn, req)
}

// PreSignedPutURL return url.Value, if not force will return an error when object exists
func (p *S3Client) PreSignedPutURL(object string, force bool) (*url.URL, error) {
	ctx := context.Background()
	object = p.fullObject(object)

	if !force {
		_, err := p.statObject(object)
		if err != nil {
			return p.cli.PresignedPutObject(ctx, p.Bucket, object, p.PresignExpiresIn)
		}
		return nil, fmt.Errorf("%s object already exists", object)
	}

	return p.cli.PresignedPutObject(ctx, p.Bucket, object, p.PresignExpiresIn)
}

func (p *S3Client) StatObject(object string) (minio.ObjectInfo, error) {
	object = p.fullObject(object)
	return p.statObject(object)
}

func (p *S3Client) DeleteObject(object string) error {
	ctx := context.Background()
	object = p.fullObject(object)

	return p.cli.RemoveObject(ctx, p.Bucket, object, minio.RemoveObjectOptions{})
}

func (p *S3Client) UploadFile(object string, file string, force bool) (minio.UploadInfo, error) {
	object = p.fullObject(object)

	info, _ := p.statObject(object)

	if info.ETag != "" && !force {
		fmt.Printf("%s Exist, skip\n", object)
		return minio.UploadInfo{}, nil
	}

	ctx := context.Background()
	n := uint(5)

	fmt.Printf("copy %s -> %s\n", file, filepath.Join(p.Bucket, object))
	return p.cli.FPutObject(ctx, p.Bucket, object, file, minio.PutObjectOptions{
		NumThreads: n,
	})
}

func (p *S3Client) UploadFolder(objectFolder string, folder string, recursive bool, force bool) error {
	finfo, err := os.Stat(folder)
	if err != nil {
		return err
	}
	if !finfo.IsDir() {
		return fmt.Errorf("not a folder")
	}

	fs, err := ioutil.ReadDir(folder)
	if err != nil {
		return nil
	}

	for _, f := range fs {
		subsrc := filepath.Join(folder, f.Name())
		subdest := filepath.Join(objectFolder, f.Name())

		if !f.Mode().IsDir() {
			_, err = p.UploadFile(subdest, subsrc, force)
			if err != nil {
				logrus.Warning(err)
			}
		}

		if f.IsDir() && recursive {
			_ = p.UploadFolder(subdest, subsrc, recursive, force)
		}
	}

	return nil
}

func (p *S3Client) fullObject(object string) string {
	return filepath.Join(p.ObjectPrefixPath, object)
}

func (p *S3Client) statObject(object string) (minio.ObjectInfo, error) {
	ctx := context.Background()
	return p.cli.StatObject(ctx, p.Bucket, object, minio.StatObjectOptions{})
}
