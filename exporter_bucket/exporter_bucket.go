package exporter_bucket

import "github.com/minio/minio-go/v7"

type ExporterBucket struct {
	minioClient *minio.Client
	bucketName  string
}
