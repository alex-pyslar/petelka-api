package handler

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"github.com/alex-pyslar/online-store/internal/models"
	"github.com/alex-pyslar/online-store/internal/service"
	"github.com/golang-jwt/jwt/v5"
	"github.com/pkg/errors"
)

// AuthHandler обрабатывает запросы авторизации с использованием JWT
type AuthHandler struct {
	userService *service.UserService
}

// NewAuthHandler создаёт новый обработчик авторизации
func NewAuthHandler(s *service.UserService) *AuthHandler {
	return &AuthHandler{userService: s}
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
	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.userService.CreateUser(r.Context(), &user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Устанавливаем статус 201 и возвращаем созданного пользователя
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
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	user, err := h.userService.GetUserByEmail(r.Context(), req.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "Invalid email or password", http.StatusUnauthorized)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := h.userService.VerifyPassword(r.Context(), user.ID, req.Password); err != nil {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	expirationTime := time.Now().Add(24 * time.Hour) // Токен действителен 24 часа
	claims := &Claims{
		UserID: user.ID,
		Email:  user.Email,
		Role:   user.Role, // Сохраняем роль в Claims
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "online-store",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(LoginResponse{Token: tokenString})
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
