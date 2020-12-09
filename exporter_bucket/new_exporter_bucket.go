package exporter_bucket

import (
	"context"
	"fmt"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/lifecycle"
)

func NewExporterBucket(
	ctx context.Context,
	minioClient *minio.Client,
	bucketName,
	region string,
) (
	st ExporterBucket,
	err error,
) {
	err = minioClient.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{
		Region: region,
	})
	if err != nil &&
		err.Error() != "Your previous request to create the named bucket succeeded and you already own it." {
		return
	}

	err = minioClient.SetBucketPolicy(ctx, bucketName, fmt.Sprintf(`{
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

	err = minioClient.SetBucketLifecycle(ctx, bucketName, &lifecycle.Configuration{
		Rules: []lifecycle.Rule{{
			ID:     "Remove expired files",
			Status: "Enabled",
			Expiration: lifecycle.Expiration{
				Days: 2,
			},
		}, {
			ID:     "Remove expired multipart upload",
			Status: "Enabled",
			AbortIncompleteMultipartUpload: lifecycle.AbortIncompleteMultipartUpload{
				DaysAfterInitiation: 1,
			},
		}},
	})
	if err != nil {
		return
	}

	st = ExporterBucket{
		minioClient: minioClient,
		bucketName:  bucketName,
	}
	return
}
