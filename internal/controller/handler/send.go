package handler

import (
	"avito-shop/internal/model"
	"avito-shop/internal/service"
	"encoding/json"
	"net/http"
)

type SendCoinHandler struct {
	coinService *service.CoinService
}

func NewSendCoinHandler(coinService *service.CoinService) *SendCoinHandler {
	return &SendCoinHandler{coinService: coinService}
}

func (h *SendCoinHandler) SendCoin(w http.ResponseWriter, r *http.Request) {
	var req entity.SendCoinRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	fromUser := r.Context().Value("username").(string)

	if err := h.coinService.SendCoin(fromUser, req.ToUser, req.Amount); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}
