package config

import (
	"context"
	"database/sql"
	"fmt"
	"os"

	"github.com/alex-pyslar/petelka-api/internal/logger"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
)

type Config struct {
	DB             *sql.DB
	Redis          *redis.Client
	MinioEndpoint  string
	MinioAccessKey string
	MinioSecretKey string
	MinioBucket    string
	MinioUseSSL    bool
}

func NewConfig(log *logger.Logger) (*Config, error) {
	// Загружаем .env файл
	if err := godotenv.Load(); err != nil {
		log.Warningf("Failed to load .env file: %v", err)
	}

	// --- PostgreSQL ---
	dbConnStr := os.Getenv("DATABASE_URL")
	if dbConnStr == "" {
		log.Fatal("DATABASE_URL is not set")
	}
	db, err := sql.Open("postgres", dbConnStr)
	if err != nil {
		log.Errorf("Failed to open PostgreSQL connection: %v", err)
		return nil, err
	}
	if err := db.Ping(); err != nil {
		log.Errorf("Failed to ping PostgreSQL: %v", err)
		return nil, err
	}
	log.Info("Connected to PostgreSQL successfully")

	// --- Redis ---
	redisAddr := os.Getenv("REDIS_URL")
	if redisAddr == "" {
		log.Fatal("REDIS_URL is not set")
	}
	redisClient := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: "",
		DB:       0,
	})
	if _, err := redisClient.Ping(context.Background()).Result(); err != nil {
		log.Errorf("Failed to ping Redis: %v", err)
		return nil, err
	}
	log.Info("Connected to Redis successfully")

	// --- MinIO ---
	minioEndpoint := os.Getenv("MINIO_ENDPOINT")
	minioAccessKey := os.Getenv("MINIO_ACCESS_KEY")
	minioSecretKey := os.Getenv("MINIO_SECRET_KEY")
	minioBucket := os.Getenv("MINIO_BUCKET")
	minioUseSSL := os.Getenv("MINIO_USE_SSL") == "true"

	if minioEndpoint == "" || minioAccessKey == "" || minioSecretKey == "" || minioBucket == "" {
		return nil, fmt.Errorf("MinIO config missing: MINIO_ENDPOINT, MINIO_ACCESS_KEY, MINIO_SECRET_KEY, MINIO_BUCKET")
	}

	log.Info("MinIO configuration loaded")

	return &Config{
		DB:             db,
		Redis:          redisClient,
		MinioEndpoint:  minioEndpoint,
		MinioAccessKey: minioAccessKey,
		MinioSecretKey: minioSecretKey,
		MinioBucket:    minioBucket,
		MinioUseSSL:    minioUseSSL,
	}, nil
}
