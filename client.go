package confs3

import (
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/mohae/deepcopy"
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

func New(akid string, akey string, endpoint string, ssl bool) *S3Client {
	return &S3Client{
		AccessID:  akid,
		AccessKey: akey,
		Endpoint:  endpoint,
		SSL:       ssl,
	}
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

func (p *S3Client) SetRegion(region string) *S3Client {
	p.Region = region
	return p
}

func (p *S3Client) SetBucket(bucket string) *S3Client {
	p.Bucket = bucket
	return p
}

func (p *S3Client) SetExpiresIn(second int) *S3Client {
	if second == 0 {
		second = 300
	}
	p.ExpiresIn = time.Duration(second) * time.Second
	return p
}

// Fork create a deepcopy s3client pointer
func (p *S3Client) Fork() *S3Client {
	return deepcopy.Copy(p).(*S3Client)
}
