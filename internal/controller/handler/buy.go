package handler

import (
	"avito-shop/internal/usecase"
	"github.com/gorilla/mux"
	"net/http"
)

type BuyHandler struct {
	buyService *usecase.BuyService
}

func NewBuyHandler(buyService *usecase.BuyService) *BuyHandler {
	return &BuyHandler{buyService: buyService}
}

func (h *BuyHandler) BuyItem(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	item := vars["item"]

	username := r.Context().Value("username").(string)

	if err := h.buyService.BuyItem(username, item); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}
