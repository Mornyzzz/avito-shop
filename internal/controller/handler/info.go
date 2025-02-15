package handler

import (
	"avito-shop/internal/entity"
	"avito-shop/internal/usecase"
	"encoding/json"
	"net/http"
)

type InfoHandler struct {
	infoService *usecase.InfoService
}

func NewInfoHandler(infoService *usecase.InfoService) *InfoHandler {
	return &InfoHandler{infoService: infoService}
}

func (h *InfoHandler) GetInfo(w http.ResponseWriter, r *http.Request) {
	username := r.Context().Value("username").(string)

	info, err := h.infoService.GetInfo(username)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(info)
}
