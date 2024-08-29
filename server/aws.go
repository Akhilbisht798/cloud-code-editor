package main

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3Client struct {
	Client *s3.Client
}

func (s3Client S3Client) GetPresignedUrls(bucketName string, prefix string) ([]string, error) {
	var urls []string
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
			url, err := presignerClient.getObject(bucketName, *obj.Key, int64(5))
			if err != nil {
				panic(err.Error())
			}
			urls = append(urls, url.URL)
		}
	}
	return urls, nil
}

func getS3ClientDevelopment() S3Client {
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

func createBucket() {
	client := getS3ClientDevelopment()
	cparams := &s3.CreateBucketInput{
		Bucket: aws.String("project"),
	}
	_, err := client.Client.CreateBucket(context.Background(), cparams)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
}

// func getS3Client() *s3.Client {
// 	cfg, err := config.LoadDefaultConfig(context.TODO())
// 	if err != nil {
// 		fmt.Println("Error: loading the aws config")
// 		log.Fatal("err")
// 	}
// 	client := s3.NewFromConfig(cfg)
// 	return client
// }

// func retriveProject(project *Project) {

// }

// func storeProject(project *Project) {
// }
