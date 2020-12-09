package minio

import (
	"context"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"strconv"
	"time"
)

func NewClient(
	ctx context.Context,
	accessKeyID,
	secretAccessKey,
	endpoint,
	secure string,
) (
	client *minio.Client,
	err error,
) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	sec, err := strconv.ParseBool(secure)
	if err != nil {
		return
	}

	client, err = minio.New(endpoint, &minio.Options{
		Creds: credentials.NewStaticV4(
			accessKeyID,
			secretAccessKey,
			"",
		),
		Secure: sec,
	})
	if err != nil {
		return
	}

	// ping
	_, err = client.ListBuckets(ctx)
	return
}
