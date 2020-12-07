package store

import "github.com/minio/minio-go/v7"

func NewStore(minioClient *minio.Client, bucketName string) Store {
	return Store{
		minioClient: minioClient,
		bucketName:  bucketName,
	}
}
