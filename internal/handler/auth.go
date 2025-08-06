package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/alex-pyslar/online-store/internal/logger"
	"github.com/alex-pyslar/online-store/internal/models"
	"github.com/alex-pyslar/online-store/internal/service"
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

// AuthHandler обрабатывает запросы авторизации с использованием JWT
type AuthHandler struct {
	userService *service.UserService
	log         *logger.Logger
}

// NewAuthHandler создаёт новый обработчик авторизации
func NewAuthHandler(s *service.UserService, log *logger.Logger) *AuthHandler {
	return &AuthHandler{userService: s, log: log}
}

// Register godoc
// @Summary Регистрация нового пользователя
// @Description Создаёт нового пользователя в системе
// @Tags auth
// @Accept json
// @Produce json
// @Param user body models.User true "Данные пользователя для регистрации"
// @Success 201 {object} models.User "Пользователь успешно создан"
// @Failure 400 {string} string "Неверный формат запроса или пользователь с таким email уже существует"
// @Failure 500 {string} string "Внутренняя ошибка сервера"
// @Router /auth/register [post]
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	requestID := uuid.New().String()
	ctx := context.WithValue(r.Context(), "request_id", requestID)

	h.log.Infof("Register attempt from IP: %s, request_id: %s", r.RemoteAddr, requestID)

	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		h.log.Errorf("Error decoding register request body, request_id: %s: %v", requestID, err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.userService.CreateUser(ctx, &user); err != nil {
		h.log.Errorf("Failed to register user, request_id: %s: %v", requestID, err)
		http.Error(w, "Failed to register user", http.StatusInternalServerError)
		return
	}

	h.log.Infof("User registered successfully: %s, request_id: %s", user.Email, requestID)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

// Login godoc
// @Summary Вход пользователя
// @Description Аутентификация пользователя и выдача JWT токена
// @Tags auth
// @Accept json
// @Produce json
// @Param credentials body LoginRequest true "Учетные данные пользователя"
// @Success 200 {object} LoginResponse "Успешная аутентификация, возвращает токен"
// @Failure 400 {string} string "Неверный запрос или отсутствуют учетные данные"
// @Failure 401 {string} string "Неверный email или пароль"
// @Failure 500 {string} string "Внутренняя ошибка сервера"
// @Router /auth/login [post]
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	requestID := uuid.New().String()
	ctx := context.WithValue(r.Context(), "request_id", requestID)

	h.log.Infof("Login attempt from IP: %s, request_id: %s", r.RemoteAddr, requestID)

	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.log.Errorf("Error decoding login request body, request_id: %s: %v", requestID, err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	user, err := h.userService.GetUserByEmail(ctx, req.Email)
	if err != nil {
		h.log.Errorf("Failed to get user by email %s, request_id: %s: %v", req.Email, requestID, err)
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	if err := h.userService.VerifyPassword(ctx, user.ID, req.Password); err != nil {
		h.log.Errorf("Password verification failed for user %s, request_id: %s: %v", user.Email, requestID, err)
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	expirationTime := time.Now().Add(24 * time.Hour) // Токен действителен 24 часа
	claims := &Claims{
		UserID: user.ID,
		Email:  user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "online-store",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		h.log.Errorf("Failed to sign token for user %s, request_id: %s: %v", user.Email, requestID, err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	h.log.Infof("User %s logged in successfully, request_id: %s", user.Email, requestID)
	json.NewEncoder(w).Encode(LoginResponse{Token: tokenString})
}

// JWTMiddleware godoc
// @Summary Проверка JWT токена
// @Description Middleware для проверки и валидации JWT токена в заголовке Authorization
// @Security ApiKeyAuth
// @Failure 401 {string} string "Неавторизованный доступ"
// @Failure 500 {string} string "Внутренняя ошибка сервера"
func (h *AuthHandler) JWTMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := uuid.New().String()
		ctx := context.WithValue(r.Context(), "request_id", requestID)
		r = r.WithContext(ctx)

		w.Header().Set("Access-Control-Allow-Origin", "https://petelka.shop")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			h.log.Errorf("Authorization header is required, request_id: %s", requestID)
			http.Error(w, "Authorization header is required", http.StatusUnauthorized)
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			h.log.Errorf("Invalid Authorization header format: %s, request_id: %s", authHeader, requestID)
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
			h.log.Errorf("Failed to parse JWT token, request_id: %s: %v", requestID, err)
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		if !token.Valid {
			h.log.Errorf("Invalid token, request_id: %s", requestID)
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		// Проверка Issuer и других полей
		if claims.Issuer != "online-store" {
			h.log.Errorf("Invalid issuer in token, request_id: %s", requestID)
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		h.log.Infof("JWT token valid for user ID: %d, email: %s, request_id: %s", claims.UserID, claims.Email, requestID)
		ctx = context.WithValue(ctx, "userID", claims.UserID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// LoginRequest представляет структуру запроса для входа
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// LoginResponse представляет структуру ответа для входа
type LoginResponse struct {
	Token string `json:"token"`
}
