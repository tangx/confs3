package confs3

import (
	"fmt"
	"testing"
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

	s := `tools/hello-go`
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

	u, err = s3.SetExpiresIn(600).PreSignedPutURL(s, false)
	if err != nil {
		panic(err)
	}
	fmt.Println(u)

}
