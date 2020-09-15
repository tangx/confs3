package confs3

import (
	"fmt"
	"testing"
)

func TestMain(t *testing.T) {
	s3 := &S3Client{
		AccessID:  "",
		AccessKey: "",
		SSL:       true,
		Bucket:    "",
		Endpoint:  "",
	}

	s3.Init()
	err := s3.Login()
	if err != nil {
		panic(err)
	}

	s := `tools/hello-go`
	u, err := s3.PreSignedGetURL(s)
	if err != nil {
		panic(err)
	}
	fmt.Println(u)

	u, err = s3.PreSignedPutURL(s, false)
	if err != nil {
		panic(err)
	}
	fmt.Println(u)
}
