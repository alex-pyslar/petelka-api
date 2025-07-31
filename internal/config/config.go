package config

import (
	"context"
	"database/sql"
	"github.com/redis/go-redis/v9"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type Config struct {
	DB          *sql.DB
	Redis       *redis.Client
	OAuthConfig *oauth2.Config
}

func NewConfig() (*Config, error) {
	// Загрузка переменных окружения
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	// Подключение к PostgreSQL
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	if err := db.Ping(); err != nil {
		return nil, err
	}

	// Подключение к Redis
	redisClient := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_URL"),
		Password: "", // Установить, если нужен пароль
		DB:       0,  // Использовать базу по умолчанию
	})
	if err := redisClient.Ping(context.Background()).Err(); err != nil {
		return nil, err
	}

	// Конфигурация OAuth
	oauthConfig := &oauth2.Config{
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		RedirectURL:  "http://localhost:8080/auth/google/callback",
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email", "https://www.googleapis.com/auth/userinfo.profile"},
		Endpoint:     google.Endpoint,
	}

	// Логирование успешного подключения с текущей датой и временем
	currentTime := time.Now()
	log.Printf("Config initialized successfully at %s", currentTime.Format("2006-01-02 15:04:05 MST"))

	return &Config{
		DB:          db,
		Redis:       redisClient,
		OAuthConfig: oauthConfig,
	}, nil
}
