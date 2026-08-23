package service

import (
	"crypto-server/internal/models"
	"crypto-server/internal/repository"
	"errors"
	"strings"
	"time"
)

type CryptoService struct {
	cryptoRepo *repository.CryptoRepository
}

func NewCryptoService(cryptoRepo *repository.CryptoRepository) *CryptoService {
	return &CryptoService{cryptoRepo: cryptoRepo}
}

func (s *CryptoService) AddCrypto(symbol string) (*models.Crypto, error) {
	if exists := s.cryptoRepo.Exists(symbol); exists {
		return nil, errors.New("crypto already exists")
	}

	crypto := models.Crypto{
		Symbol:       strings.ToUpper(symbol),
		Name:         "CRYPTO: " + symbol,
		CurrentPrice: 0.0,
		LastUpdated:  time.Now(),
	}

	s.cryptoRepo.SaveCrypto(&crypto)
	return &crypto, nil
}

func (s *CryptoService) ListAllCryptos() []*models.Crypto {
	return s.cryptoRepo.GetAll()
}

func (s *CryptoService) GetCryptoBySymbol(symbol string) (*models.Crypto, error) {
	if exists := s.cryptoRepo.Exists(symbol); !exists {
		return nil, errors.New("Crypto " + symbol + " not found")
	}

	crypto := s.cryptoRepo.FindCrypto(symbol)

	return crypto, nil
}

func (s *CryptoService) DeleteCrypto(symbol string) error {
	if exists := s.cryptoRepo.Exists(symbol); !exists {
		return errors.New("No such crypto yet")
	}

	s.cryptoRepo.DeleteCrypto(symbol)
	return nil
}
