package storage

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func ListObjects(ctx context.Context, client *s3.Client, bucket string) ([]string, error) {
	resp, err := client.ListObjectsV2(ctx, &s3.ListObjectsV2Input{
		Bucket: &bucket,
	})
	if err != nil {
		return nil, err
	}

	var keys []string
	for _, item := range resp.Contents {
		keys = append(keys, *item.Key)
	}
	return keys, nil
}
