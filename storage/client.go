package storage

import (
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func NewSpacesClient() *s3.Client {
	cfg := aws.Config{
		Credentials: credentials.NewStaticCredentialsProvider(
			os.Getenv("DO_SPACES_KEY"),
			os.Getenv("DO_SPACES_SECRET"),
			"",
		),
		Region: "us-east-1", // required but ignored by DO
		EndpointResolver: aws.EndpointResolverFunc(func(service, region string) (aws.Endpoint, error) {
			return aws.Endpoint{
				URL:           os.Getenv("DO_SPACES_ENDPOINT"),
				SigningRegion: os.Getenv("DO_SPACES_REGION"),
			}, nil
		}),
	}
	return s3.NewFromConfig(cfg)
}
