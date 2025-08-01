package config

import (
	"context"
	"database/sql"
	"github.com/alex-pyslar/online-store/internal/logger"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
)

type Config struct {
	DB    *sql.DB
	Redis *redis.Client
}

func NewConfig(log *logger.Logger) (*Config, error) {
	// Загружаем .env файл
	if err := godotenv.Load(); err != nil {
		log.Warningf("Failed to load .env file: %v", err)
	}

	// PostgreSQL
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

	// Redis
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

	return &Config{DB: db, Redis: redisClient}, nil
}
