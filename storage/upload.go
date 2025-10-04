package storage

import (
	"bytes"
	"context"
	"fmt"
	"mime/multipart"
	"os"

	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func UploadObject(ctx context.Context, client *s3.Client, bucket, key string, file multipart.File, size int64, contentType string) (string, error) {
	buffer := make([]byte, size)
	_, err := file.Read(buffer)
	if err != nil {
		return "", err
	}

	_, err = client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      &bucket,
		Key:         &key,
		Body:        bytes.NewReader(buffer),
		ContentType: &contentType,
		ACL:         "public-read",
	})
	if err != nil {
		return "", err
	}

	url := fmt.Sprintf("%s/%s/%s", os.Getenv("DO_SPACES_ENDPOINT"), bucket, key)
	return url, nil
}
