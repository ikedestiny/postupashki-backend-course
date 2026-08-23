package service

import (
	"crypto-server/internal/models"
	"crypto-server/internal/repository"
	"errors"
	"fmt"
	"strings"
	"time"
)

type CryptoService struct {
	repo  *repository.CryptoRepository
	gecko *CoinGeckoClient
}

func NewCryptoService(repo *repository.CryptoRepository, gecko *CoinGeckoClient) *CryptoService {
	return &CryptoService{
		repo:  repo,
		gecko: gecko,
	}
}

func (s *CryptoService) AddCrypto(symbol string) (*models.Crypto, error) {
	// Convert to uppercase
	symbol = strings.ToUpper(symbol)

	// Check if crypto already exists
	if s.repo.Exists(symbol) {
		return nil, errors.New("crypto already exists")
	}

	// Get CoinGecko ID
	coinID, err := s.gecko.GetCoinID(symbol)
	if err != nil {
		return nil, fmt.Errorf("failed to get CoinGecko ID: %w", err)
	}

	// Fetch real price
	price, err := s.gecko.FetchPrice(coinID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch price: %w", err)
	}

	// Create crypto with real data
	crypto := &models.Crypto{
		Symbol:       symbol,
		Name:         symbol, // You could fetch the name from CoinGecko too
		CurrentPrice: price,
		LastUpdated:  time.Now(),
	}

	// Save to repository
	s.repo.SaveCrypto(crypto)
	return crypto, nil
}

func (s *CryptoService) ListAllCryptos() []*models.Crypto {
	return s.repo.GetAll()
}

func (s *CryptoService) GetCryptoBySymbol(symbol string) (*models.Crypto, error) {
	if exists := s.repo.Exists(symbol); !exists {
		return nil, errors.New("Crypto " + symbol + " not found")
	}

	crypto := s.repo.FindCrypto(symbol)

	return crypto, nil
}

func (s *CryptoService) DeleteCrypto(symbol string) error {
	if exists := s.repo.Exists(symbol); !exists {
		return errors.New("No such crypto yet")
	}

	s.repo.DeleteCrypto(symbol)
	return nil
}
