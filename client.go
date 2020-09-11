package confs3

import "github.com/minio/minio-go/v7"

type S3Client struct {
	Endpoint  string `env:"endpoint"`
	Bucket    string `env:"bucket"`
	SSL       bool   `env:"ssl"`
	AccessID  string `env:"access_id"`
	AccessKey string `env:"access_key"`
	ExpiresIn int    `env:"expires_in" default:"300" summary:"expired time for presign url (second)"`

	*minio.Client
}

func NewClient() *S3Client {
	return &S3Client{}
}
