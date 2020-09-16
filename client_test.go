package confs3

import (
	"testing"
)

func TestMain(t *testing.T) {
	s3 := &S3Client{
		AccessID:  "AKID123456",
		AccessKey: "AKEY123456",
		SSL:       false,
		Bucket:    "",
		Endpoint:  "127.0.0.1:9000",
	}

	s3.Init()
	err := s3.Login()
	if err != nil {
		panic(err)
	}

	err = s3.CreateBucket("confs3")
	if err != nil {
		panic(err)
	}
	// s := `tools/hello-go`
	// u, err := s3.PreSignedGetURL(s)
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println(u)

	// u, err = s3.PreSignedPutURL(s, false)
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println(u)

}
