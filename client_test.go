package confs3

import (
	"fmt"
	"log"
	"testing"

	. "github.com/onsi/gomega"
)

const (
	bucket = "confs3"
)

var (
	s3 = &S3Client{
		AccessID:  "AKID123456",
		AccessKey: "AKEY123456",
		Endpoint:  "127.0.0.1:9000",
		SSL:       false,
		Bucket:    bucket,
	}
)

func Test_PreSign(t *testing.T) {
	s3.Init()
	u, err := s3.PreSignedGetURL("gorm-demo.tgz")
	if err != nil {
		panic(err)
	}
	fmt.Println(u)
}

func TestMain(t *testing.T) {

	s3.Init()

	err := s3.CreateBucket(bucket)
	t.Run("CreateBucket", func(t *testing.T) {
		NewWithT(t).Expect(err).To(BeNil())
	})

	s3.SetBucket(bucket).SetExpiresIn(30)

	s := `avantar.jpg`

	u, err := s3.PreSignedGetURL(s)
	t.Run("PreSignedGetURL", func(t *testing.T) {
		NewWithT(t).Expect(err).To(BeNil())
		NewWithT(t).Expect(u.Hostname()).NotTo(Equal(s3.Endpoint))
		log.Print(u)
	})

	u, err = s3.SetExpiresIn(600).PreSignedPutURL(s, true)
	t.Run("PreSignedPutURL_SetExpiresIn", func(t *testing.T) {
		NewWithT(t).Expect(err).To(BeNil())
		NewWithT(t).Expect(u.Hostname()).NotTo(Equal(s3.Endpoint))
		log.Print(u)
	})

}

func TestUpload(t *testing.T) {
	src := `./`
	dest := `path/2/test-folder2/`

	s3.Init()

	_ = s3.CreateBucket(bucket)

	s3.SetBucket(bucket)

	err := s3.UploadFolder(dest, src, true, true)
	t.Run("UploadFolder", func(t *testing.T) {
		NewWithT(t).Expect(err).To(BeNil())
	})
}
