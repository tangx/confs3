package confs3

import (
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type S3Client struct {
	Endpoint  string        `env:"endpoint"`
	Bucket    string        `env:"bucket"`
	Region    string        `env:"region"`
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

func (p *S3Client) SetRegion(region string) {
	p.Region = region
}

func (p *S3Client) SetBucket(bucket string) {
	p.Bucket = bucket
}
