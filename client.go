package confs3

import (
	"strings"
	"sync"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/sirupsen/logrus"
)

type S3Client struct {
	Endpoint         string        `env:""`
	Bucket           string        `env:""`
	Region           string        `env:""`
	SSL              bool          `env:""`
	AccessID         string        `env:""`
	AccessKey        string        `env:""`
	ObjectPrefixPath string        `env:""`
	PresignExpiresIn time.Duration `env:"" default:"300" comment:"expired time for presign url (second)"`

	cli *minio.Client
}

var lock = sync.Mutex{}

func (p *S3Client) Init() {

	lock.Lock()
	defer lock.Unlock()

	if p.cli == nil {
		p.SetDefaults()
		p.initial()
	}
}

func (p *S3Client) initial() {

	cli, err := minio.New(p.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(p.AccessID, p.AccessKey, ""),
		Secure: p.SSL,
	})
	if err != nil {
		logrus.Fatal("s3client create failed")
		return
	}

	p.cli = cli

}

func (p *S3Client) SetDefaults() {
	if int64(p.PresignExpiresIn) == 0 {
		p.PresignExpiresIn = 300 * time.Second
	}
	p.ObjectPrefixPath = strings.Trim(p.ObjectPrefixPath, "/")
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
	p.PresignExpiresIn = time.Duration(second) * time.Second
	return p
}

/*
	功能本身没有问题， 因为要实现单例模式， 所以注释
	github.com/mohae/deepcopy
*/
// // Fork create a deepcopy s3client pointer
// func (p *S3Client) Fork() *S3Client {
// 	return deepcopy.Copy(p).(*S3Client)
// }
