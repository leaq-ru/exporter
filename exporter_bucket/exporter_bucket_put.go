package exporter_bucket

import (
	"context"
	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"os"
	"path/filepath"
	"time"
)

func (s ExporterBucket) Put(ctx context.Context, path string, deleteAfterUpload bool) (s3URL string, err error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Minute)
	defer cancel()

	if deleteAfterUpload {
		defer os.Remove(path)
	}

	ext := filepath.Ext(path)

	u, err := uuid.NewRandom()
	if err != nil {
		return
	}

	objectName := u.String() + ext

	file, err := os.Open(path)
	if err != nil {
		return
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		return
	}

	obj, err := s.minioClient.PutObject(ctx, s.bucketName, objectName, file, stat.Size(), minio.PutObjectOptions{})
	if err != nil {
		return
	}

	s3URL = "https://" + s.bucketName + ".ru/" + obj.Key
	return
}
