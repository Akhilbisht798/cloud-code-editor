package cloud

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

var ENV = os.Getenv("APP_ENV")

type S3Client struct {
	Client *s3.Client
}

func (s3Client S3Client) GetPresignedUrls(bucketName string, prefix string) (map[string]string, error) {
	vals := make(map[string]string)

	paginator := s3.NewListObjectsV2Paginator(s3Client.Client, &s3.ListObjectsV2Input{
		Bucket: aws.String(bucketName),
		Prefix: aws.String(prefix),
	})

	for paginator.HasMorePages() {
		page, err := paginator.NextPage(context.TODO())
		if err != nil {
			return nil, fmt.Errorf("failed to get page: %v", err)
		}

		for _, obj := range page.Contents {
			url, err := PresignerClient.getObject(bucketName, *obj.Key, int64(5))
			if err != nil {
				panic(err.Error())
			}
			vals[url.URL] = *obj.Key
		}
	}
	return vals, nil
}

func GetS3ClientDevelopment() S3Client {
	log.Println("Getting s3 development client")
	const defautlRegion = "us-east-1"
	staticResolver := aws.EndpointResolverFunc(func(service, region string) (aws.Endpoint, error) {
		return aws.Endpoint{
			PartitionID:       "aws",
			URL:               "http://localhost:9000",
			SigningRegion:     defautlRegion,
			HostnameImmutable: true,
		}, nil
	})

	cfg := aws.Config{
		Region:           defautlRegion,
		Credentials:      credentials.NewStaticCredentialsProvider("ROOT", "password", ""),
		EndpointResolver: staticResolver,
	}

	s3Client := s3.NewFromConfig(cfg)
	client := S3Client{
		Client: s3Client,
	}
	return client
}

func CreateBucket() {
	var client S3Client
	if ENV == "production" {
		client = GetS3Client()
	} else {
		client = GetS3ClientDevelopment()
	}
	cparams := &s3.CreateBucketInput{
		Bucket: aws.String("project"),
	}
	_, err := client.Client.CreateBucket(context.Background(), cparams)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
}
func GetS3Client() S3Client {
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithDefaultRegion("us-east-1"),
	)
	if err != nil {
		fmt.Println("Error: loading the aws config")
		log.Fatal("err")
	}
	log.Println("S3 Production client created")
	client := s3.NewFromConfig(cfg)
	c := S3Client{
		Client: client,
	}
	return c
}

func CreateNewProjectS3(key string) error {
	env := os.Getenv("APP_ENV")
	var client S3Client
	if env == "production" {
		client = GetS3Client()
	} else {
		client = GetS3ClientDevelopment()
	}
	str := "### README.md\nHello World"
	body := strings.NewReader(str)
	k := key + "README.md"

	input := s3.PutObjectInput{
		Bucket: aws.String(os.Getenv("BUCKET")),
		Key:    aws.String(k),
		Body:   body,
	}
	_, err := client.Client.PutObject(context.TODO(), &input)
	if err != nil {
		return err
	}
	return nil
}
