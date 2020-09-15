package confs3

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"path/filepath"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type S3Client struct {
	Endpoint  string        `env:"endpoint"`
	Bucket    string        `env:"bucket"`
	SSL       bool          `env:"ssl"`
	AccessID  string        `env:"access_id"`
	AccessKey string        `env:"access_key"`
	ExpiresIn time.Duration `env:"expires_in" default:"300" comment:"expired time for presign url (second)"`

	cli *minio.Client
}

func NewClient() *S3Client {
	return &S3Client{}
}

func (p *S3Client) Init() {
	p.SetDefaults()

}

func (p *S3Client) SetDefaults() {
	if int64(p.ExpiresIn) == 0 {
		p.ExpiresIn = 300 * time.Second
	}
}

func (p *S3Client) Login() (err error) {

	p.cli, err = minio.New(p.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(p.AccessID, p.AccessKey, ""),
		Secure: p.SSL,
	})
	if err != nil {
		return err
	}

	return nil
}

func (p *S3Client) PreSignedGetURL(object string) (*url.URL, error) {
	ctx := context.Background()
	req := make(url.Values)
	fname := filepath.Base(object)
	req.Set(`response-content-disposition`, fmt.Sprintf(`attachment; filename="%s"`, fname))
	return p.cli.PresignedGetObject(ctx, p.Bucket, object, p.ExpiresIn, req)
}

// PreSignedPutURL return url.Value, if not force will return an error when object exists
func (p *S3Client) PreSignedPutURL(object string, force bool) (*url.URL, error) {
	ctx := context.Background()

	if force {
		return p.cli.PresignedPutObject(ctx, p.Bucket, object, p.ExpiresIn)
	}

	obj, err := p.cli.StatObject(ctx, p.Bucket, object, minio.StatObjectOptions{})
	if err != nil {
		return p.cli.PresignedPutObject(ctx, p.Bucket, object, p.ExpiresIn)
	}
	spew.Dump(obj)

	msg := fmt.Sprintf("%s object already exists", object)
	return nil, errors.New(msg)
}
