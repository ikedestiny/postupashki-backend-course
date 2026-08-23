package repository

import (
	"hw2/internal/models"
	"log"
	"maps"
	"slices"
	"sync"
)

type CryptoRepository struct {
	mu      sync.RWMutex
	cryptos map[string]*models.Crypto
	history map[string][]models.PriceRecord
}

func NewCryptoRepository() *CryptoRepository {
	return &CryptoRepository{
		cryptos: make(map[string]*models.Crypto),
		history: make(map[string][]models.PriceRecord),
	}
}

func (r *CryptoRepository) SaveCrypto(crypto *models.Crypto) {
	r.mu.Lock()
	defer r.mu.Unlock()

	// ADD DEBUG:
	log.Printf("Saving crypto: %s", crypto.Symbol)

	r.cryptos[crypto.Symbol] = crypto
}

func (r *CryptoRepository) FindCrypto(symbol string) *models.Crypto {
	r.mu.RLock()
	defer r.mu.RUnlock()

	// ADD DEBUG:
	log.Printf("Looking for symbol: %s", symbol)
	log.Printf("Current map keys: %v", getMapKeys(r.cryptos)) // You'll need to implement this

	return r.cryptos[symbol]
}

func (r *CryptoRepository) DeleteCrypto(symbol string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.cryptos, symbol)
	delete(r.history, symbol)
}

func (r *CryptoRepository) AddHistory(symbol string, record models.PriceRecord) {
	r.mu.Lock()
	defer r.mu.Unlock()

	hist := r.history[symbol]
	if hist == nil {
		hist = []models.PriceRecord{}
	}

	hist = append(hist, record)

	//keep only last 100
	if len(hist) > 100 {
		hist = hist[len(hist)-100:] //drop first 100
	}

	r.history[symbol] = hist

}

// GetHistory returns a COPY of the history slice.
func (r *CryptoRepository) GetHistory(symbol string) []models.PriceRecord {
	r.mu.RLock()
	defer r.mu.RUnlock()
	// Return a copy so the caller can't accidentally modify the internal map
	hist := r.history[symbol]
	if hist == nil {
		return []models.PriceRecord{}
	}
	return hist
}

func (r *CryptoRepository) GetAll() []*models.Crypto {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return slices.Collect(maps.Values(r.cryptos))
}

func (r *CryptoRepository) Exists(symbol string) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()
	_, exists := r.cryptos[symbol]
	return exists
}

// Helper function for debugging
func getMapKeys(m map[string]*models.Crypto) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}
