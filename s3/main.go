package main

import (
	"crypto/md5"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

func main() {
	h := md5.New()
	content := strings.NewReader("")
	content.WriteTo(h)

	// Initialize a session in cn-north-1 that the SDK will use to load
	// credentials from the shared credentials file ~/.aws/credentials.
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("cn-northwest-1")},
	)

	// Create S3 service client
	svc := s3.New(sess)
	now := time.Now()
	resp, _ := svc.PutObjectRequest(&s3.PutObjectInput{
		Bucket: aws.String(os.Getenv("AWS_S3_BUCKET")),
		Key: aws.String(fmt.Sprintf("%d/%d/%d/%d/%d/%d",
			now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second())),
	})

	md5s := base64.StdEncoding.EncodeToString(h.Sum(nil))
	resp.HTTPRequest.Header.Set("Content-MD5", md5s)

	url, err := resp.Presign(15 * time.Minute)
	if err != nil {
		fmt.Println("error presigning request", err)
		return
	}
	fmt.Println("pressigned url:", url)
	// body 与 md5 要匹配
	req, err := http.NewRequest("PUT", url, strings.NewReader(""))
	req.Header.Set("Content-MD5", md5s)
	if err != nil {
		fmt.Println("error creating request", url)
		return
	}

	putResp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("put failed.", err)
		return
	}
	defer putResp.Body.Close()
	io.Copy(os.Stdout, putResp.Body)
}
