package store

import "github.com/minio/minio-go/v7"

type Store struct {
	minioClient *minio.Client
	bucketName  string
}
