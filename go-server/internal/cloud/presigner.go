package cloud

import (
	"context"
	"log"
	"os"
	"time"

	v4 "github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go/aws"
)

type Presigner struct {
	PresignClient *s3.PresignClient
}

var presignerClient *Presigner

func getPresignerClient() {
	env := os.Getenv("APP_ENV")
	var s3Client S3Client
	if env == "production" {
		s3Client = GetS3Client()
	} else {
		s3Client = GetS3ClientDevelopment()
	}
	client := s3.NewPresignClient(s3Client.Client)
	signer := &Presigner{
		PresignClient: client,
	}
	presignerClient = signer
}

func (presigner Presigner) putObject(
	bucketName string, objectKey string, lifetimeSecs int64) (*v4.PresignedHTTPRequest, error) {
	request, err := presigner.PresignClient.PresignPutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectKey),
	}, func(opts *s3.PresignOptions) {
		opts.Expires = time.Duration(lifetimeSecs * int64(time.Second))
	})
	if err != nil {
		log.Printf("Couldn't get a presigned request to put %v:%v. Here's why: %v\n",
			bucketName, objectKey, err)
	}
	return request, err
}

func (presigner Presigner) getObject(
	bucketName string, objecetKey string, lifetimeSecs int64) (*v4.PresignedHTTPRequest, error) {
	request, err := presigner.PresignClient.PresignGetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objecetKey),
	}, func(opts *s3.PresignOptions) {
		opts.Expires = time.Duration(lifetimeSecs * int64(time.Minute))
	})
	if err != nil {
		log.Printf("Could'nt get a Presigned request to get %v:%v, Here why: %v\n",
			bucketName, objecetKey, err)
	}
	return request, err
}
