package handler

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/alex-pyslar/online-store/internal/models"
	"github.com/alex-pyslar/online-store/internal/service"
	"golang.org/x/oauth2"
)

type AuthHandler struct {
	oauthConfig *oauth2.Config
	userService *service.UserService
}

func NewAuthHandler(oauthConfig *oauth2.Config, userService *service.UserService) *AuthHandler {
	return &AuthHandler{oauthConfig: oauthConfig, userService: userService}
}

// @Summary Initiate Google OAuth login
// @Description Redirects to Google for authentication
// @Tags auth
// @Success 302
// @Router /auth/google/login [get]
func (h *AuthHandler) HandleGoogleLogin(w http.ResponseWriter, r *http.Request) {
	url := h.oauthConfig.AuthCodeURL("state")
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

// @Summary Handle Google OAuth callback
// @Description Handles the callback from Google OAuth and creates/updates user
// @Tags auth
// @Param code query string true "Authorization code"
// @Success 200 {string} string "User authenticated"
// @Router /auth/google/callback [get]
func (h *AuthHandler) HandleGoogleCallback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	token, err := h.oauthConfig.Exchange(context.Background(), code)
	if err != nil {
		http.Error(w, "Failed to exchange token", http.StatusBadRequest)
		return
	}

	client := h.oauthConfig.Client(context.Background(), token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		http.Error(w, "Failed to get user info", http.StatusBadRequest)
		return
	}
	defer resp.Body.Close()

	var userInfo struct {
		ID    string `json:"id"`
		Email string `json:"email"`
		Name  string `json:"name"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		http.Error(w, "Failed to decode user info", http.StatusBadRequest)
		return
	}

	user := &models.User{
		Email:   userInfo.Email,
		Name:    userInfo.Name,
		OAuthID: userInfo.ID,
	}
	if err := h.userService.CreateUser(r.Context(), user); err != nil {
		http.Error(w, "Failed to save user", http.StatusInternalServerError)
		return
	}

	w.Write([]byte("User authenticated: " + userInfo.Email))
}
