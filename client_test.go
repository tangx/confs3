package confs3

import (
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
	}
)

func TestMain(t *testing.T) {

	// s3 := New("AKID123456", "", "127.0.0.1:9000", false)

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
	})

	/* fork and panic */
	// s3new := s3.Fork().SetBucket("s3conf")
	// u, err = s3new.PreSignedGetURL(s)
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println(u)

	u, err = s3.SetExpiresIn(600).PreSignedPutURL(s, true)
	t.Run("PreSignedPutURL", func(t *testing.T) {
		NewWithT(t).Expect(err).To(BeNil())
		NewWithT(t).Expect(u.Hostname()).NotTo(Equal(s3.Endpoint))
	})

	err = s3.SetBucketLifecycleExpireIn("/private", 333)
	t.Run("PreSignedPutURL", func(t *testing.T) {
		NewWithT(t).Expect(err).To(BeNil())
	})

	// err = s3.SetObjectPrefixExpireAt("/public", "2020-09-30 00:00:00")
	// if err != nil {
	// 	panic(err)
	// }

	_, err = s3.GetBucketLifecycle()
	if err != nil {
		panic(err)
	}
	// spew.Dump(info)
}

func TestUpload(t *testing.T) {
	src := `/tmp/test-folder`
	dest := `path/2/test-folder2/`

	s3.Init()

	_ = s3.CreateBucket(bucket)

	s3.SetBucket(bucket)

	err := s3.UploadFolder(dest, src, true, true)
	t.Run("UploadFolder", func(t *testing.T) {
		NewWithT(t).Expect(err).To(BeNil())
	})
}
