package handler

import (
	"context"
	"net/http"
	"os"
	"strings"

	"github.com/alex-pyslar/online-store/internal/logger"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// Ключ для подписи JWT
var jwtKey = []byte(os.Getenv("JWT_SECRET"))

// Claims - структура для токена JWT
type Claims struct {
	UserID int    `json:"user_id"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

// ContextKey используется для передачи данных через контекст
type ContextKey string

const (
	UserIDKey    ContextKey = "userID"
	RequestIDKey ContextKey = "requestID"
)

// JWTMiddleware проверяет и валидирует JWT токен в заголовке Authorization.
func JWTMiddleware(log *logger.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestID := uuid.New().String()
			ctx := context.WithValue(r.Context(), RequestIDKey, requestID)
			r = r.WithContext(ctx)

			log.Infof("Request received: %s %s, request_id: %s", r.Method, r.URL.Path, requestID)

			// Настройка CORS
			w.Header().Set("Access-Control-Allow-Origin", "https://petelka.shop")
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, DELETE")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}

			// Проверка заголовка авторизации
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				log.Errorf("Authorization header is required, request_id: %s", requestID)
				http.Error(w, "Authorization header is required", http.StatusUnauthorized)
				return
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				log.Errorf("Invalid Authorization header format: %s, request_id: %s", authHeader, requestID)
				http.Error(w, "Invalid Authorization header format", http.StatusUnauthorized)
				return
			}

			tokenString := parts[1]
			claims := &Claims{}

			token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, jwt.ErrSignatureInvalid
				}
				return jwtKey, nil
			})

			if err != nil {
				log.Errorf("Failed to parse JWT token, request_id: %s: %v", requestID, err)
				http.Error(w, "Invalid token", http.StatusUnauthorized)
				return
			}

			if !token.Valid {
				log.Errorf("Invalid token, request_id: %s", requestID)
				http.Error(w, "Invalid token", http.StatusUnauthorized)
				return
			}

			// Проверка Issuer
			if claims.Issuer != "online-store" {
				log.Errorf("Invalid issuer in token, request_id: %s", requestID)
				http.Error(w, "Invalid token", http.StatusUnauthorized)
				return
			}

			log.Infof("JWT token valid for user ID: %d, request_id: %s", claims.UserID, requestID)

			// Добавляем ID пользователя в контекст для последующих хендлеров
			ctx = context.WithValue(ctx, UserIDKey, claims.UserID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
