package handler

import (
	"crypto-server/internal/models"
	"crypto-server/internal/service"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type CryptoHandler struct {
	cryptoService *service.CryptoService
}

type CreateRequest struct {
	Symbol string `json:"symbol"`
}

type CreateResponse struct {
	Crypto *models.Crypto `json:"crypto"`
}

func NewCryptohandler(cryptoService *service.CryptoService) *CryptoHandler {
	return &CryptoHandler{cryptoService: cryptoService}
}

func (h *CryptoHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req CreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error": "invalid JSON"}`, http.StatusBadRequest)
		return
	}

	if req.Symbol == "" {
		http.Error(w, "symbol cannot be empty", http.StatusBadRequest)
		return
	}

	crypto, err := h.cryptoService.AddCrypto(req.Symbol)

	if err != nil {
		if err.Error() == "crypto already exists" {
			http.Error(w, "already exists", http.StatusConflict)
		}
		http.Error(w, "something went wrong while trying to add crypto", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(CreateResponse{Crypto: crypto})
}

func (h *CryptoHandler) ListAll(w http.ResponseWriter, r *http.Request) {
	cryptos := h.cryptoService.ListAllCryptos()

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(cryptos)
}

// GetBySymbol handles GET /api/crypto/{symbol}
func (h *CryptoHandler) GetBySymbol(w http.ResponseWriter, r *http.Request) {
	symbol := chi.URLParam(r, "symbol")

	crypto, err := h.cryptoService.GetCryptoBySymbol(symbol)

	if err != nil {
		http.Error(w, "no such crypto bud", http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(crypto)
}

func (h *CryptoHandler) Delete(w http.ResponseWriter, r *http.Request) {
	symbol := chi.URLParam(r, "symbol")

	err := h.cryptoService.DeleteCrypto(symbol)

	if err != nil {
		http.Error(w, "no such crypto", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("{}")
}
