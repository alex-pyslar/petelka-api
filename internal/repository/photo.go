package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"path/filepath"
	"time"

	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/redis/go-redis/v9"
)

type PhotoRepository struct {
	client     *minio.Client
	bucketName string
	redis      *redis.Client
}

func NewPhotoRepository(
	minioEndpoint, accessKey, secretKey, bucket string,
	useSSL bool, redis *redis.Client,
) (*PhotoRepository, error) {

	client, err := minio.New(minioEndpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		return nil, err
	}

	ctx := context.Background()
	exists, err := client.BucketExists(ctx, bucket)
	if err != nil {
		return nil, err
	}
	if !exists {
		if err = client.MakeBucket(ctx, bucket, minio.MakeBucketOptions{}); err != nil {
			return nil, err
		}
	}

	return &PhotoRepository{
		client:     client,
		bucketName: bucket,
		redis:      redis,
	}, nil
}

func (r *PhotoRepository) Upload(ctx context.Context, file io.Reader, size int64, filename, contentType string) (string, string, error) {
	ext := filepath.Ext(filename)
	objectName := uuid.New().String() + ext

	_, err := r.client.PutObject(ctx, r.bucketName, objectName, file, size, minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		return "", "", err
	}

	url, err := r.client.PresignedGetObject(ctx, r.bucketName, objectName, 7*24*time.Hour, nil)
	if err != nil {
		return "", "", err
	}

	presigned := url.String()

	// Инвалидация кэша
	cacheKey := fmt.Sprintf("photo_url:%s", objectName)
	r.redis.Del(ctx, cacheKey)

	return objectName, presigned, nil
}

func (r *PhotoRepository) GetPresignedURL(ctx context.Context, objectName string) (string, error) {
	cacheKey := fmt.Sprintf("photo_url:%s", objectName)

	// Попробуем из кэша
	if cached, err := r.redis.Get(ctx, cacheKey).Result(); err == nil {
		var url string
		if json.Unmarshal([]byte(cached), &url) == nil {
			return url, nil
		}
	}

	// Генерируем
	url, err := r.client.PresignedGetObject(ctx, r.bucketName, objectName, 7*24*time.Hour, nil)
	if err != nil {
		return "", err
	}
	presigned := url.String()

	// Кэшируем
	data, _ := json.Marshal(presigned)
	r.redis.Set(ctx, cacheKey, data, 7*24*time.Hour)

	return presigned, nil
}
