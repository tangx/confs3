package confs3

import (
	"fmt"
	"testing"

	"github.com/davecgh/go-spew/spew"
)

const (
	bucket = "confs3"
)

func TestMain(t *testing.T) {

	s3 := New("AKID123456", "AKEY123456", "127.0.0.1:9000", false)

	s3.Init()
	err := s3.Login()
	if err != nil {
		panic(err)
	}

	err = s3.CreateBucket(bucket)
	if err != nil {
		panic(err)
	}

	s3.SetBucket(bucket).SetExpiresIn(30)

	s := `avantar.jpg`
	u, err := s3.PreSignedGetURL(s)
	if err != nil {
		panic(err)
	}
	fmt.Println(u)

	/* fork and panic */
	// s3new := s3.Fork().SetBucket("s3conf")
	// u, err = s3new.PreSignedGetURL(s)
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println(u)

	u, err = s3.SetExpiresIn(600).PreSignedPutURL(s, true)
	if err != nil {
		panic(err)
	}
	fmt.Println(u)

	err = s3.SetBucketLifecycleExpireIn("/private", 333)
	if err != nil {
		panic(err)
	}

	// err = s3.SetObjectPrefixExpireAt("/public", "2020-09-30 00:00:00")
	// if err != nil {
	// 	panic(err)
	// }

	info, err := s3.GetBucketLifecycle()
	if err != nil {
		panic(err)
	}
	spew.Dump(info)
}

func TestUpload(t *testing.T) {
	src := `/tmp/test-folder`
	dest := `path/2/test-folder2/`

	s3 := New("AKID123456", "AKEY123456", "127.0.0.1:9000", false)

	s3.Init()
	err := s3.Login()
	if err != nil {
		panic(err)
	}

	err = s3.CreateBucket(bucket)
	if err != nil {
		panic(err)
	}
	s3.SetBucket(bucket)

	_ = s3.UploadFolder(dest, src, true)
}
