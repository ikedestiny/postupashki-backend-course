package service

import (
	"crypto-server/internal/models"
	"errors"
	"sync"
	"time"
)

type ScheduleService struct {
	mu            sync.RWMutex
	config        *models.ScheduleConfig
	cryptoService *CryptoService
	stopChan      chan struct{}
	ticker        *time.Ticker
}

func NewScheduleService(cryptoService *CryptoService) *ScheduleService {
	return &ScheduleService{
		config: &models.ScheduleConfig{
			Enabled:         true,
			IntervalSeconds: 30, // Default: 30 seconds
			LastUpdate:      time.Now(),
			NextUpdate:      time.Now().Add(30 * time.Second),
		},
		cryptoService: cryptoService,
		stopChan:      make(chan struct{}),
	}
}

// Start begins the background scheduler
func (s *ScheduleService) Start() {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.ticker != nil {
		// Already running, stop old ticker
		s.stopTickerLocked()
	}

	if s.config.Enabled {
		s.startTickerLocked()
	}
}

func (s *ScheduleService) stopTickerLocked() {
	if s.ticker != nil {
		s.ticker.Stop()
		s.ticker = nil
	}
}

func (s *ScheduleService) startTickerLocked() {
	interval := time.Duration(s.config.IntervalSeconds) * time.Second
	s.ticker = time.NewTicker(interval)

	go func() {
		for {
			select {
			case <-s.ticker.C:
				s.performUpdate()
			case <-s.stopChan:
				return
			}
		}
	}()

	// Update next update time
	s.config.NextUpdate = time.Now().Add(interval)
}

func (s *ScheduleService) performUpdate() {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.config.Enabled {
		return
	}

	// Get all cryptos
	cryptos := s.cryptoService.ListAllCryptos()
	if len(cryptos) == 0 {
		return
	}

	// Refresh each crypto
	for _, crypto := range cryptos {
		_, err := s.cryptoService.RefreshPrice(crypto.Symbol)
		if err != nil {
			// Log error but continue with other cryptos
			continue
		}
	}

	s.config.LastUpdate = time.Now()
	s.config.NextUpdate = time.Now().Add(time.Duration(s.config.IntervalSeconds) * time.Second)
}

// TriggerManualUpdate forces an immediate update
func (s *ScheduleService) TriggerManualUpdate() (int, error) {
	cryptos := s.cryptoService.ListAllCryptos()
	if len(cryptos) == 0 {
		return 0, nil
	}

	updatedCount := 0
	for _, crypto := range cryptos {
		_, err := s.cryptoService.RefreshPrice(crypto.Symbol)
		if err == nil {
			updatedCount++
		}
	}

	// Update last update time
	s.mu.Lock()
	s.config.LastUpdate = time.Now()
	if s.config.Enabled {
		s.config.NextUpdate = time.Now().Add(time.Duration(s.config.IntervalSeconds) * time.Second)
	}
	s.mu.Unlock()

	return updatedCount, nil
}

// GetConfig returns the current schedule configuration
func (s *ScheduleService) GetConfig() *models.ScheduleConfig {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// Return a copy to avoid external modifications
	config := *s.config
	return &config
}

// UpdateConfig updates the schedule configuration
func (s *ScheduleService) UpdateConfig(updated *models.ScheduleRequest) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Update enabled if provided
	if updated.Enabled != nil {
		s.config.Enabled = *updated.Enabled
	}

	// Update interval if provided (with validation)
	if updated.IntervalSeconds != nil {
		interval := *updated.IntervalSeconds
		if interval < 10 || interval > 3600 {
			return errors.New("interval must be between 10 and 3600 seconds")
		}
		s.config.IntervalSeconds = interval
	}

	// Restart the ticker with new settings
	s.stopTickerLocked()
	if s.config.Enabled {
		s.startTickerLocked()
	}

	return nil
}

// Stop gracefully stops the scheduler
func (s *ScheduleService) Stop() {
	s.mu.Lock()
	defer s.mu.Unlock()
	close(s.stopChan)
	s.stopTickerLocked()
}
