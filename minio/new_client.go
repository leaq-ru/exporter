package minio

import (
	"context"
	"fmt"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/minio/minio-go/v7/pkg/lifecycle"
	"strconv"
	"time"
)

func NewClient(
	ctx context.Context,
	accessKeyID,
	secretAccessKey,
	endpoint,
	region,
	bucketName,
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
	if err != nil {
		return
	}

	err = client.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{
		Region: region,
	})
	if err != nil {
		return
	}

	err = client.SetBucketPolicy(ctx, bucketName, fmt.Sprintf(`{
		"Version": "2012-10-17",
		"Statement": [{
			"Sid": "PublicRead",
			"Effect": "Allow",
			"Principal": "*",
			"Action": ["s3:GetObject"],
			"Resource": ["arn:aws:s3:::%s/*"]
		}]
	}`, bucketName))
	if err != nil && err.Error() != "200 OK" {
		return
	}

	err = client.SetBucketLifecycle(ctx, bucketName, &lifecycle.Configuration{
		Rules: []lifecycle.Rule{{
			ID:     "Remove expired files",
			Status: "Enabled",
			Expiration: lifecycle.Expiration{
				Days: 1,
			},
		}, {
			ID:     "Remove expired multipart upload",
			Status: "Enabled",
			AbortIncompleteMultipartUpload: lifecycle.AbortIncompleteMultipartUpload{
				DaysAfterInitiation: 1,
			},
		}},
	})
	return
}
