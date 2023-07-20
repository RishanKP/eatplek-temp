package services

import (
	 "github.com/aws/aws-sdk-go/aws"
	 "github.com/aws/aws-sdk-go/aws/credentials"
	 "github.com/aws/aws-sdk-go/aws/session"
	 "os"
	 "fmt"
	 "github.com/aws/aws-sdk-go/service/s3/s3manager"
	 "github.com/gin-gonic/gin"
)

func ConnectAws() *session.Session {
	AccessKeyID := os.Getenv("AWS_ACCESS_KEY")
	SecretAccessKey := os.Getenv("AWS_SECRET_ACCESS_KEY")
	MyRegion := os.Getenv("AWS_REGION")
	
	sess, err := session.NewSession(
	 &aws.Config{
	  Region: aws.String(MyRegion),
	  Credentials: credentials.NewStaticCredentials(
	   AccessKeyID,
	   SecretAccessKey,
	   "", // a token will be created when the session it's used.
	  ),
	 })
	
	 if err != nil {
	  panic(err)
	 }
	
	 return sess
}

func UploadFile(c *gin.Context) (string, error) {

	sess := ConnectAws()
	uploader := s3manager.NewUploader(sess)
	bucket := os.Getenv("AWS_BUCKET")

	file, header, err := c.Request.FormFile("file")

	if err != nil {
		return "", err
	}

	key := "files/"
	key += header.Filename

	filetype := header.Header["Content-Type"][0]

	up, err := uploader.Upload(&s3manager.UploadInput{
		Bucket:             aws.String(bucket),
		ACL:                aws.String("public-read"),
		Key:                aws.String(key),
		ContentType:        aws.String(filetype),
		ContentDisposition: aws.String("inline"),
		Body:               file,
	})

	if err != nil {
		fmt.Println(err)
		return "", err
	}

	return up.Location, nil
}
