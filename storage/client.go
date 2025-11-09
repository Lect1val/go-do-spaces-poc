package storage

import (
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func NewSpacesClient() *s3.Client {
	endpoint := os.Getenv("DO_SPACES_ENDPOINT")
	cfg := aws.Config{
		Credentials: credentials.NewStaticCredentialsProvider(
			os.Getenv("DO_SPACES_KEY"),
			os.Getenv("DO_SPACES_SECRET"),
			"",
		),
		Region:       os.Getenv("DO_SPACES_REGION"),
		BaseEndpoint: &endpoint,
	}
	return s3.NewFromConfig(cfg)
}
