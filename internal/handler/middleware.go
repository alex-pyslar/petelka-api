package handler

import (
	"context"
	"net/http"
	"os"
	"strings"

	"github.com/alex-pyslar/petelka-api/internal/logger"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// Ключ для подписи JWT
var jwtKey = []byte(os.Getenv("JWT_SECRET"))

// Claims - структура для токена JWT
type Claims struct {
	UserID int    `json:"user_id"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

// ContextKey используется для передачи данных через контекст
type ContextKey string

const (
	UserIDKey   ContextKey = "userID"
	UserRoleKey ContextKey = "userRole"
)

// AuthMiddleware - универсальный middleware для проверки авторизации
func AuthMiddleware(log *logger.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestID := uuid.New().String()
			ctx := context.WithValue(r.Context(), "requestID", requestID)
			r = r.WithContext(ctx)

			// Пропускаем OPTIONS-запросы
			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}

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

			if err != nil || !token.Valid {
				log.Errorf("Invalid token, request_id: %s: %v", requestID, err)
				http.Error(w, "Invalid token", http.StatusUnauthorized)
				return
			}

			// Добавляем ID и роль пользователя в контекст
			ctx = context.WithValue(ctx, UserIDKey, claims.UserID)
			ctx = context.WithValue(ctx, UserRoleKey, claims.Role) // Сохраняем роль

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// AdminMiddleware проверяет, является ли пользователь администратором.
func AdminMiddleware(log *logger.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestID := r.Context().Value("requestID").(string)
			role, ok := r.Context().Value(UserRoleKey).(string)
			if !ok || role != "admin" {
				log.Warningf("Access denied: non-admin user trying to access admin route, request_id: %s", requestID)
				http.Error(w, "Forbidden", http.StatusForbidden)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

// CorsMiddleware устанавливает заголовки CORS для всех запросов.
func CorsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "https://petelka.shop")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		next.ServeHTTP(w, r)
	})
}
