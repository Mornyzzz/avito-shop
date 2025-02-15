package handler

import (
	"avito-shop/internal/entity"
	"avito-shop/internal/usecase"
	"encoding/json"
	"net/http"
)

type AuthHandler struct {
	authService *usecase.AuthService
}

func NewAuthHandler(authService *usecase.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

func (h *AuthHandler) Authenticate(w http.ResponseWriter, r *http.Request) {
	var req entity.AuthRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	token, err := h.authService.Authenticate(req.Username, req.Password)
	if err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	response := entity.AuthResponse{Token: token}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
